<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { api } from "$lib/api";
  import { fmtDate, textPreview } from "$lib/mailutils";
  import type { Mail, Mailbox } from "$lib/types";

  let mailboxId = $derived(Number(page.params.id));

  let mailbox = $state<Mailbox | null>(null);
  let mail = $state<Mail[] | undefined>(undefined);
  let syncing = $state(false);
  let error = $state("");
  let notice = $state("");

  async function load() {
    error = "";
    try {
      const boxes = await api.listMailboxes();
      mailbox = boxes.find((b) => b.id === mailboxId) ?? null;
      mail = await api.listMail(mailboxId);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load mail";
      mail = undefined;
    }
  }

  onMount(load);

  async function handleSync() {
    syncing = true;
    error = "";
    notice = "";
    try {
      await api.syncMailbox(mailboxId);
      notice = "Mailbox synced";
      await load();
    } catch (e) {
      error = e instanceof Error ? e.message : "Sync failed";
    } finally {
      syncing = false;
    }
  }
</script>

<div class="mailbox-wrap">
  <div class="page-head">
    <div class="crumbs">
      <a class="crumb" href="/">box</a>
      <span class="sep">/</span>
      <span class="crumb current">{mailbox?.username ?? `Mailbox ${mailboxId}`}</span>
    </div>
    <button class="btn" onclick={handleSync} disabled={syncing}>
      {syncing ? "Syncing…" : "↻ Sync"}
    </button>
  </div>

  {#if error}
    <div class="error-banner">{error}</div>
  {/if}
  {#if notice}
    <div class="success-banner">{notice}</div>
  {/if}

  {#if mail === undefined && !error}
    <p class="dim pad">Loading mail…</p>
  {:else if mail !== undefined && mail.length === 0}
    <div class="empty">
      <p><strong>Inbox is empty</strong></p>
      <p class="dim">Click sync to fetch new messages from the server.</p>
    </div>
  {:else if mail !== undefined}
    <ul class="mail-list">
      {#each mail as m (m.id)}
        <li>
          <a class="mail-row" href={`/mailbox/${mailboxId}/message/${m.id}`}>
            <span class="from">{m.display_name || m.from_address}</span>
            <span class="mid">
              <span class="subject">{m.subject || "(no subject)"}</span>
              <span class="preview">{textPreview(m.body_parts)}</span>
            </span>
            <span class="date">{fmtDate(m.recieved_at)}</span>
          </a>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .mailbox-wrap {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 860px;
    margin: 0 auto;
    padding: 1.5rem;
  }

  .page-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .crumbs {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-weight: 600;
    font-size: 0.9rem;
    color: #767676;
    min-width: 0;
  }

  .crumb {
    color: var(--primary-color);
    text-decoration: none;
  }

  .crumb.current {
    color: inherit;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .sep {
    opacity: 0.5;
  }

  .btn {
    background: linear-gradient(180deg, var(--primary-color), var(--primary-color-dark));
    border: 1px solid #308cb0;
    color: white;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
    padding: 0.5rem 0.9rem;
    border-radius: 10px;
    font-size: 0.9rem;
    cursor: pointer;
    white-space: nowrap;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: wait;
  }

  .error-banner,
  .success-banner {
    border-radius: 8px;
    padding: 0.6rem 0.9rem;
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }

  .error-banner {
    background: #fdecea;
    color: #a33;
    border: 1px solid #f5c6c0;
  }

  .success-banner {
    background: #eafaf1;
    color: #1e7a43;
    border: 1px solid #b9e4cb;
  }

  .dim {
    color: #767676;
  }

  .pad {
    padding: 1rem 0;
  }

  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: #333;
  }

  .mail-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .mail-row {
    display: grid;
    grid-template-columns: 180px 1fr auto;
    align-items: baseline;
    gap: 1rem;
    padding: 0.8rem 1rem;
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 10px;
    color: #222;
    text-decoration: none;
    transition: border-color 0.12s ease, background 0.12s ease;
  }

  .mail-row:hover {
    border-color: var(--primary-color);
    background: #f7fbfd;
  }

  .from {
    font-weight: 650;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .mid {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .subject {
    color: #222;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .preview {
    color: #959595;
    font-size: 0.82rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .date {
    color: #959595;
    font-size: 0.8rem;
    white-space: nowrap;
  }

  @media (max-width: 640px) {
    .mail-row {
      grid-template-columns: 1fr;
      gap: 0.25rem;
    }
    .date {
      font-size: 0.75rem;
    }
  }
</style>
