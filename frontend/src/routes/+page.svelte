<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Mailbox } from '$lib/api';

	let mailboxes: Mailbox[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let syncingId: number | null = $state(null);

	let showAdd = $state(false);
	let form = $state({ server_hostport: '', username: '', password: '', primary_inbox: 'INBOX' });
	let adding = $state(false);
	let addError = $state('');
	let notice = $state('');

	onMount(async () => {
		await loadMailboxes();
	});

	async function loadMailboxes() {
		loading = true;
		error = '';
		try {
			mailboxes = await api.listMailboxes();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load mailboxes';
		} finally {
			loading = false;
		}
	}

	async function handleSync(id: number) {
		syncingId = id;
		error = '';
		notice = '';
		try {
			await api.syncMailbox(id);
			notice = 'Mailbox synced';
			await loadMailboxes();
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
			await loadMailboxes();
		} catch (e) {
			addError = e instanceof Error ? e.message : 'Failed to add mailbox';
		} finally {
			adding = false;
		}
	}
</script>

<div class="dashboard">
	<div class="page-head">
		<div>
			<h1 class="page-title">Mailboxes</h1>
			<p class="page-sub">Add an email account to start syncing your inbox.</p>
		</div>
		<button class="btn btn-primary" onclick={() => (showAdd = !showAdd)}>
			{showAdd ? 'Cancel' : '+ Add mailbox'}
		</button>
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

	{#if loading}
		<div class="spinner-center"><div class="spinner" style="width:28px;height:28px"></div></div>
	{:else if mailboxes.length === 0}
		<div class="empty">
			<p style="margin:0 0 6px;font-size:17px">No mailboxes yet</p>
			<p style="margin:0">Click "Add mailbox" above to get started.</p>
		</div>
	{:else}
		<ul class="mailbox-list">
			{#each mailboxes as mailbox (mailbox.id)}
				<li class="card mailbox-row">
					<a class="mailbox-main" href={`/mailbox/${mailbox.id}`}>
						<span class="mailbox-avatar">
							{mailbox.username.slice(0, 1).toUpperCase()}
						</span>
						<span class="mailbox-info">
							<span class="mailbox-username">{mailbox.username}</span>
							<span class="mailbox-meta">
								{mailbox.server_hostport} · {mailbox.primary_inbox}
							</span>
						</span>
					</a>
					<div class="mailbox-actions">
						<button
							class="btn btn-sm btn-ghost"
							onclick={() => handleSync(mailbox.id)}
							disabled={syncingId === mailbox.id}
						>
							{#if syncingId === mailbox.id}
								<span class="spinner"></span> Syncing…
							{:else}
								⟳ Sync
							{/if}
						</button>
						<a class="btn btn-sm" href={`/mailbox/${mailbox.id}`}>Open</a>
					</div>
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

	.add-form {
		padding: 20px;
		margin-bottom: 24px;
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

	.mailbox-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.mailbox-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 14px 16px;
	}

	.mailbox-main {
		display: flex;
		align-items: center;
		gap: 14px;
		min-width: 0;
		color: var(--text);
		text-decoration: none !important;
	}

	.mailbox-avatar {
		width: 42px;
		height: 42px;
		flex-shrink: 0;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--accent), #8a5bef);
		display: grid;
		place-items: center;
		font-size: 18px;
		font-weight: 700;
		color: #fff;
	}

	.mailbox-info {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.mailbox-username {
		font-weight: 600;
		font-size: 16px;
	}

	.mailbox-meta {
		color: var(--text-dim);
		font-size: 13px;
	}

	.mailbox-actions {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-shrink: 0;
	}

	@media (max-width: 640px) {
		.form-grid {
			grid-template-columns: 1fr;
		}
		.mailbox-row {
			flex-direction: column;
			align-items: stretch;
		}
	}
</style>
