<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api";
  import type { Mailbox } from "$lib/types";

  let mailboxes = $state<Mailbox[] | undefined>(undefined);
  let loading = $state(true);
  let error = $state("");
  let notice = $state("");
  let syncingID = $state<number | undefined>(undefined);
  let removingID = $state<number | undefined>(undefined);

  let showAdd = $state(false);
  let adding = $state(false);
  let addError = $state("");
  let form = $state({
    server_hostport: "",
    username: "",
    password: "",
    primary_inbox: "INBOX",
  });

  async function load() {
    loading = true;
    error = "";
    try {
      mailboxes = await api.listMailboxes();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load mailboxes";
    } finally {
      loading = false;
    }
  }

  onMount(load);

  async function handleAdd() {
    addError = "";
    adding = true;
    try {
      await api.createMailbox({
        server_hostport: form.server_hostport.trim(),
        username: form.username.trim(),
        password: form.password,
        primary_inbox: form.primary_inbox.trim() || "INBOX",
      });
      form = { server_hostport: "", username: "", password: "", primary_inbox: "INBOX" };
      showAdd = false;
      await load();
    } catch (e) {
      addError = e instanceof Error ? e.message : "Failed to add mailbox";
    } finally {
      adding = false;
    }
  }

  async function handleSync(id: number) {
    syncingID = id;
    error = "";
    notice = "";
    try {
      await api.syncMailbox(id);
      notice = "Mailbox synced";
    } catch (e) {
      error = e instanceof Error ? e.message : "Sync failed";
    } finally {
      syncingID = undefined;
    }
  }

  async function handleRemove(id: number) {
    if (!confirm("Remove this mailbox?")) return;
    removingID = id;
    error = "";
    notice = "";
    try {
      await api.deleteMailbox(id);
      notice = "Mailbox removed";
      await load();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to remove mailbox";
    } finally {
      removingID = undefined;
    }
  }
</script>

<div class="manage-wrap">
  <div class="page-head">
    <div>
      <h1 class="title">Mailboxes</h1>
      <p class="sub">Add an email account to start syncing your inbox.</p>
    </div>
    <button class="btn" onclick={() => (showAdd = !showAdd)}>
      {showAdd ? "Cancel" : "+ Add mailbox"}
    </button>
  </div>

  {#if error}
    <div class="error-banner">{error}</div>
  {/if}
  {#if notice}
    <div class="success-banner">{notice}</div>
  {/if}

  {#if showAdd}
    <form class="card add-form" onsubmit={(e) => { e.preventDefault(); handleAdd(); }}>
      <h2 class="form-title">New mailbox</h2>
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
        <div class="error-banner add-error">{addError}</div>
      {/if}
      <div class="form-actions">
        <button class="btn" type="submit" disabled={adding}>{adding ? "Adding…" : "Add mailbox"}</button>
      </div>
    </form>
  {/if}

  {#if loading}
    <p class="dim">Loading…</p>
  {:else if mailboxes && mailboxes.length === 0}
    <div class="empty">
      <p><strong>No mailboxes yet</strong></p>
      <p class="dim">Click "Add mailbox" above to get started.</p>
    </div>
  {:else if mailboxes}
    <ul class="mailbox-list">
      {#each mailboxes as mailbox (mailbox.id)}
        <li class="card row">
          <a class="main" href={`/mailbox/${mailbox.id}`}>
            <span class="avatar">{mailbox.username.slice(0, 1).toUpperCase()}</span>
            <span class="info">
              <span class="username">{mailbox.username}</span>
              <span class="meta">{mailbox.server_hostport} · {mailbox.primary_inbox}</span>
            </span>
          </a>
          <div class="actions">
            <button class="btn btn-sm" onclick={() => handleSync(mailbox.id)} disabled={syncingID === mailbox.id}>
              {syncingID === mailbox.id ? "Syncing…" : "↻ Sync"}
            </button>
            <button
              class="btn btn-sm btn-danger"
              onclick={() => handleRemove(mailbox.id)}
              disabled={removingID === mailbox.id}
            >
              {removingID === mailbox.id ? "Removing…" : "Remove"}
            </button>
            <a class="btn btn-sm" href={`/mailbox/${mailbox.id}`}>Open</a>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .manage-wrap {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 760px;
    margin: 0 auto;
    padding: 1.5rem;
  }

  .page-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .title {
    font-size: 1.5rem;
    margin: 0;
  }

  .sub {
    color: #767676;
    margin: 0.25rem 0 0;
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

  .btn-sm {
    padding: 0.35rem 0.7rem;
    font-size: 0.82rem;
  }

  .btn-danger {
    background: linear-gradient(180deg, #e05c5c, #c74040);
    border-color: #b03838;
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

  .card {
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 12px;
    background: #fff;
  }

  .form-title {
    font-size: 1.1rem;
    margin: 0 0 1rem;
  }

  .add-form {
    padding: 1.25rem;
    margin-bottom: 1.5rem;
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
  }

  .field label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #444;
  }

  .field input {
    padding: 0.5rem 0.6rem;
    border: 1px solid rgba(0, 0, 0, 0.15);
    border-radius: 8px;
    font-size: 0.9rem;
  }

  .add-error {
    margin-top: 1rem;
  }

  .form-actions {
    margin-top: 1.25rem;
    display: flex;
    justify-content: flex-end;
  }

  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: #333;
  }

  .dim {
    color: #767676;
  }

  .mailbox-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.9rem 1rem;
  }

  .main {
    display: flex;
    align-items: center;
    gap: 0.9rem;
    min-width: 0;
    color: #222;
    text-decoration: none;
  }

  .avatar {
    width: 42px;
    height: 42px;
    flex-shrink: 0;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--primary-color), #8a5bef);
    display: grid;
    place-items: center;
    font-size: 18px;
    font-weight: 700;
    color: #fff;
  }

  .info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .username {
    font-weight: 600;
  }

  .meta {
    color: #959595;
    font-size: 0.8rem;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }

  @media (max-width: 640px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
    .row {
      flex-direction: column;
      align-items: stretch;
    }
    .actions {
      flex-wrap: wrap;
    }
  }
</style>
