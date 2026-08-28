<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Mailbox, type MessageWithMailbox } from '$lib/api';

	let mailboxes: Mailbox[] = $state([]);
	let messages: MessageWithMailbox[] = $state([]);
	let loading = $state(true);
	let syncingAll = $state(false);
	let syncingId: number | null = $state(null);
	let error = $state('');
	let notice = $state('');

	let showAdd = $state(false);
	let form = $state({ server_hostport: '', username: '', password: '', primary_inbox: 'INBOX' });
	let adding = $state(false);
	let addError = $state('');

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		error = '';
		notice = '';
		try {
			const { mailboxes: boxes, messages: all, failed } = await api.listAllMessages();
			mailboxes = boxes;
			messages = all;
			if (failed > 0) error = `Failed to load messages from ${failed} mailbox(es)`;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load messages';
		} finally {
			loading = false;
		}
	}

	async function handleSyncAll() {
		syncingAll = true;
		error = '';
		notice = '';
		try {
			const results = await Promise.allSettled(mailboxes.map((b) => api.syncMailbox(b.id)));
			const failed = results.filter((r) => r.status === 'rejected').length;
			if (failed > 0) {
				error = `Failed to sync ${failed} mailbox(es)`;
			} else {
				notice = results.length > 0 ? 'All mailboxes synced' : 'No mailboxes to sync';
			}
			await load();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Sync failed';
		} finally {
			syncingAll = false;
		}
	}

	async function handleSync(id: number) {
		syncingId = id;
		error = '';
		notice = '';
		try {
			await api.syncMailbox(id);
			notice = 'Mailbox synced';
			await load();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Sync failed';
		} finally {
			syncingId = null;
		}
	}

	async function handleAdd() {
		addError = '';
		adding = true;
		try {
			await api.createMailbox({
				server_hostport: form.server_hostport.trim(),
				username: form.username.trim(),
				password: form.password,
				primary_inbox: form.primary_inbox.trim() || 'INBOX'
			});
			form = { server_hostport: '', username: '', password: '', primary_inbox: 'INBOX' };
			showAdd = false;
			await load();
		} catch (e) {
			addError = e instanceof Error ? e.message : 'Failed to add mailbox';
		} finally {
			adding = false;
		}
	}

	function fmtDate(iso: string): string {
		const d = new Date(iso);
		if (isNaN(d.getTime())) return iso;
		return d.toLocaleString(undefined, {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function preview(from: string): string {
		const m = from.match(/^(.*?)(?:\s*<[^>]*>)?$/);
		const name = (m ? m[1] : from).trim().replace(/^"|"$/g, '');
		return name || from;
	}

	function inboxName(mailbox: Mailbox): string {
		return mailbox.username || mailbox.server_hostport;
	}
</script>

<div class="inbox">
	<div class="page-head">
		<div>
			<h1 class="page-title">Inbox</h1>
			<p class="page-sub">All mail from {mailboxes.length} mailbox{mailboxes.length === 1 ? '' : 'es'}.</p>
		</div>
		<div class="head-actions">
			<button class="btn" onclick={handleSyncAll} disabled={syncingAll || mailboxes.length === 0}>
				{#if syncingAll}<span class="spinner"></span> Syncing…{:else}⟳ Sync all{/if}
			</button>
			<button
				class="btn btn-primary"
				onclick={() => {
					addError = '';
					showAdd = !showAdd;
				}}
			>
				{showAdd ? 'Cancel' : '+ Add mailbox'}
			</button>
		</div>
	</div>

	{#if error}
		<div class="error-banner" style="margin-bottom:16px">{error}</div>
	{/if}

	{#if notice}
		<div class="success-banner" style="margin-bottom:16px">{notice}</div>
	{/if}

	{#if showAdd}
		<form class="card add-form" onsubmit={(e) => { e.preventDefault(); handleAdd(); }}>
			<h2 style="margin:0 0 16px">New mailbox</h2>
			<div class="form-grid">
				<div class="field">
					<label for="host">IMAP host:port</label>
					<input id="host" bind:value={form.server_hostport} placeholder="imap.example.com:993" required />
				</div>
				<div class="field">
					<label for="username">Username / email</label>
					<input id="username" bind:value={form.username} placeholder="you@example.com" autocomplete="username" required />
				</div>
				<div class="field">
					<label for="password">Password</label>
					<input id="password" bind:value={form.password} type="password" autocomplete="current-password" required />
				</div>
				<div class="field">
					<label for="inbox">Primary inbox</label>
					<input id="inbox" bind:value={form.primary_inbox} placeholder="INBOX" />
				</div>
			</div>
			{#if addError}
				<div class="error-banner" style="margin-top:16px">{addError}</div>
			{/if}
			<div class="form-actions">
				<button class="btn btn-primary" type="submit" disabled={adding}>
					{#if adding}<span class="spinner"></span>{/if}
					{adding ? 'Adding…' : 'Add mailbox'}
				</button>
			</div>
		</form>
	{/if}

	{#if mailboxes.length > 0}
		<ul class="accounts">
			{#each mailboxes as mailbox (mailbox.id)}
				<li>
					<a class="account-chip" href={`/mailbox/${mailbox.id}`} title={mailbox.server_hostport}>
						<span class="inbox-dot"></span>
						{inboxName(mailbox)}
					</a>
					<button
						class="btn btn-sm btn-ghost"
						onclick={() => handleSync(mailbox.id)}
						disabled={syncingId === mailbox.id}
					>
						{#if syncingId === mailbox.id}
							<span class="spinner" style="width:14px;height:14px"></span> Syncing…
						{:else}
							⟳ Sync
						{/if}
					</button>
				</li>
			{/each}
		</ul>
	{/if}

	{#if loading}
		<div class="spinner-center"><div class="spinner" style="width:28px;height:28px"></div></div>
	{:else if mailboxes.length === 0}
		<div class="empty">
			<p style="margin:0 0 6px;font-size:17px">No mailboxes yet</p>
			<p style="margin:0">Click "+ Add mailbox" above to get started.</p>
		</div>
	{:else if messages.length === 0}
		<div class="empty">
			<p style="margin:0 0 6px;font-size:17px">No mail yet</p>
			<p style="margin:0">Click "Sync all" to fetch new messages from your mailboxes.</p>
		</div>
	{:else}
		<ul class="msg-list">
			{#each messages as msg (msg.id)}
				<li>
					<a class="card msg-row" href={`/mailbox/${msg.mailbox.id}/message/${msg.id}`}>
						<span class="msg-from">{preview(msg.from_address)}</span>
						<span class="msg-subject">{msg.subject || '(no subject)'}</span>
						<span class="msg-inbox" title={`Delivered to ${inboxName(msg.mailbox)} · ${msg.mailbox.primary_inbox}`}>
							<span class="inbox-dot"></span>
							{inboxName(msg.mailbox)}
						</span>
						<span class="msg-date">{fmtDate(msg.received_at)}</span>
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.page-head {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 4px;
		position: relative;
		z-index: 5;
	}

	.head-actions {
		display: flex;
		align-items: center;
		gap: 10px;
		flex-shrink: 0;
	}

	.add-form {
		padding: 20px;
		margin: 20px 0 8px;
	}

	.form-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 16px;
	}

	.form-actions {
		margin-top: 20px;
		display: flex;
		gap: 12px;
		justify-content: flex-end;
	}

	.accounts {
		list-style: none;
		margin: 20px 0 20px;
		padding: 0;
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 8px;
	}

	.accounts li {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		background: var(--bg-elev-2);
		border: 1px solid var(--border);
		border-radius: 999px;
		padding: 3px 6px 3px 3px;
	}

	.account-chip {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 4px 6px 4px 10px;
		border-radius: 999px;
		color: var(--text);
		font-size: 13px;
		font-weight: 600;
		text-decoration: none !important;
	}

	.account-chip:hover {
		color: var(--accent-strong);
	}

	.inbox-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--accent);
		flex-shrink: 0;
	}

	.msg-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.msg-row {
		display: grid;
		grid-template-columns: 170px minmax(0, 1fr) auto auto;
		align-items: baseline;
		gap: 16px;
		padding: 13px 16px;
		color: var(--text);
		text-decoration: none !important;
		transition: border-color 0.12s ease, background 0.12s ease;
	}

	.msg-row:hover {
		border-color: var(--accent);
		background: var(--bg-elev-2);
	}

	.msg-from {
		font-weight: 650;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.msg-subject {
		color: var(--text-dim);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.msg-inbox {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		max-width: 220px;
		color: var(--text-dim);
		font-size: 13px;
		font-weight: 600;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.msg-date {
		color: var(--text-dim);
		font-size: 13px;
		white-space: nowrap;
	}

	@media (max-width: 860px) {
		.msg-row {
			grid-template-columns: 1fr;
			gap: 4px;
		}
		.msg-inbox {
			max-width: none;
		}
		.msg-date {
			font-size: 12px;
		}
	}

	@media (max-width: 640px) {
		.form-grid {
			grid-template-columns: 1fr;
		}
		.page-head {
			flex-direction: column;
		}
	}
</style>