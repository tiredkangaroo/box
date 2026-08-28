<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api, type Email, type Mailbox } from '$lib/api';

	let mailboxId = $derived(Number(page.params.id));
	let emailId = $derived(Number(page.params.message_id));

	let mailbox: Mailbox | null = $state(null);
	let email: Email | null = $state(null);
	let loading = $state(true);
	let error = $state('');

	function fmtDate(iso: string): string {
		const d = new Date(iso);
		if (isNaN(d.getTime())) return iso;
		return d.toLocaleString(undefined, {
			weekday: 'short',
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function cleanName(from: string): string {
		return from.replace(/^<|>$/g, '');
	}

	onMount(async () => {
		loading = true;
		try {
			const boxes = await api.listMailboxes();
			mailbox = boxes.find((b) => b.id === mailboxId) ?? null;
			email = await api.getMessage(mailboxId, emailId);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load message';
		} finally {
			loading = false;
		}
	});
</script>

<div>
	<div class="crumbs">
		<a href="/">All mail</a>
		<span class="sep">/</span>
		<a href={`/mailbox/${mailboxId}`}>{mailbox?.username ?? `Mailbox ${mailboxId}`}</a>
	</div>

	{#if loading}
		<div class="spinner-center"><div class="spinner" style="width:28px;height:28px"></div></div>
	{:else if error}
		<div class="error-banner" style="margin-top:16px">{error}</div>
	{:else if email}
		<article class="card message">
			<div class="message-head">
				<div class="avatar">{cleanName(email.from_address).slice(0, 1).toUpperCase()}</div>
				<div class="head-text">
					<div class="from">{cleanName(email.from_address)}</div>
					<div class="subject">{email.subject || '(no subject)'}</div>
					<div class="date">
						{fmtDate(email.received_at)}
						{#if mailbox}
							<span class="inbox-tag">→ {mailbox.username}{mailbox.primary_inbox ? ` · ${mailbox.primary_inbox}` : ''}</span>
						{/if}
					</div>
				</div>
			</div>
			<hr class="rule" />
			<div class="body">
				{#if email.body}
					{email.body}
				{:else}
					<span class="dim">No body.</span>
				{/if}
			</div>
		</article>
		<div class="back">
			<a class="btn btn-ghost" href="/">← Back to all mail</a>
		</div>
	{/if}
</div>

<style>
	.crumbs {
		display: flex;
		align-items: center;
		gap: 6px;
		color: var(--text-dim);
		font-weight: 600;
		font-size: 14px;
		margin-bottom: 20px;
	}

	.sep {
		opacity: 0.5;
	}

	.message {
		padding: 24px;
	}

	.message-head {
		display: flex;
		gap: 16px;
		align-items: center;
	}

	.avatar {
		width: 48px;
		height: 48px;
		flex-shrink: 0;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--accent), #8a5bef);
		display: grid;
		place-items: center;
		font-size: 20px;
		font-weight: 700;
		color: #fff;
	}

	.head-text {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.from {
		font-size: 15px;
		font-weight: 650;
	}

	.subject {
		font-size: 20px;
		font-weight: 700;
	}

	.date {
		color: var(--text-dim);
		font-size: 13px;
	}

	.inbox-tag {
		margin-left: 8px;
		color: var(--accent-strong);
	}

	.rule {
		border: none;
		border-top: 1px solid var(--border);
		margin: 20px 0;
	}

	.body {
		white-space: pre-wrap;
		word-break: break-word;
		line-height: 1.6;
		font-size: 15px;
	}

	.dim {
		color: var(--text-dim);
	}

	.back {
		margin-top: 20px;
	}
</style>
