<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api, type EmailSummary, type Mailbox } from '$lib/api';

	let mailboxId = $derived(Number(page.params.id));

	let mailbox: Mailbox | null = $state(null);
	let messages: EmailSummary[] = $state([]);
	let loading = $state(true);
	let syncing = $state(false);
	let error = $state('');
	let notice = $state('');

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
		// strip surrounding angle brackets / quotes, take everything before '@'
		const m = from.match(/^(.*?)(?:\s*<[^>]*>)?$/);
		const name = (m ? m[1] : from).trim().replace(/^"|"$/g, '');
		return name || from;
	}

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		error = '';
		try {
			const boxes = await api.listMailboxes();
			mailbox = boxes.find((b) => b.id === mailboxId) ?? null;
			messages = await api.listMessages(mailboxId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load messages';
		} finally {
			loading = false;
		}
	}

	async function handleSync() {
		syncing = true;
		error = '';
		notice = '';
		try {
			await api.syncMailbox(mailboxId);
			notice = 'Mailbox synced';
			await load();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Sync failed';
		} finally {
			syncing = false;
		}
	}
</script>

<div>
	<div class="page-head">
		<div class="crumbs">
			<a href="/">Mailboxes</a>
			<span class="sep">/</span>
			<span>{mailbox?.username ?? `Mailbox ${mailboxId}`}</span>
		</div>
		<button class="btn" onclick={handleSync} disabled={syncing}>
			{#if syncing}<span class="spinner"></span> Syncing…{:else}⟳ Sync{/if}
		</button>
	</div>

	{#if error}
		<div class="error-banner" style="margin-bottom:16px">{error}</div>
	{/if}

	{#if notice}
		<div class="success-banner" style="margin-bottom:16px">{notice}</div>
	{/if}

	{#if loading}
		<div class="spinner-center"><div class="spinner" style="width:28px;height:28px"></div></div>
	{:else if messages.length === 0}
		<div class="empty">
			<p style="margin:0 0 6px;font-size:17px">Inbox is empty</p>
			<p style="margin:0">Click sync to fetch new messages from the server.</p>
		</div>
	{:else}
		<ul class="msg-list">
			{#each messages as msg (msg.id)}
				<li>
					<a class="card msg-row" href={`/mailbox/${mailboxId}/message/${msg.id}`}>
						<span class="msg-from">{preview(msg.from_address)}</span>
						<span class="msg-subject">{msg.subject || '(no subject)'}</span>
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
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		margin-bottom: 20px;
	}

	.crumbs {
		display: flex;
		align-items: center;
		gap: 6px;
		color: var(--text-dim);
		font-weight: 600;
		font-size: 14px;
	}

	.sep {
		opacity: 0.5;
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
		grid-template-columns: 200px 1fr auto;
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

	.msg-date {
		color: var(--text-dim);
		font-size: 13px;
		white-space: nowrap;
	}

	@media (max-width: 720px) {
		.msg-row {
			grid-template-columns: 1fr;
			gap: 4px;
		}
		.msg-date {
			font-size: 12px;
		}
	}
</style>
