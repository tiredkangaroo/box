<script lang="ts">
  import "./layout.css";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { api } from "$lib/api";
  import type { Mailbox } from "$lib/types";
  import favicon from "$lib/assets/favicon.svg";

  let { children } = $props();

  let mailboxes = $state<Mailbox[] | undefined>(undefined);
  let syncingID = $state<number | undefined>(undefined);
  let loadError = $state("");

  let activeID = $derived.by(() => {
    const id = page.params.id;
    return id !== undefined ? Number(id) : undefined;
  });

  async function load() {
    loadError = "";
    try {
      mailboxes = await api.listMailboxes();
    } catch (e) {
      loadError = e instanceof Error ? e.message : "Failed to load mailboxes";
    }
  }

  onMount(load);

  async function handleSync(e: Event, id: number) {
    e.stopPropagation();
    syncingID = id;
    try {
      await api.syncMailbox(id);
      if (getActiveMailbox(id)) {
        goto(`/mailbox/${id}`, { invalidateAll: true });
      }
    } finally {
      syncingID = undefined;
    }
  }

  function getActiveMailbox(id: number) {
    return activeID === id;
  }

  function select(id: number) {
    goto(`/mailbox/${id}`);
  }
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<nav>
  <div>
    <h1 class="text-center">box</h1>
    <h2 class="nav-heading">Mailboxes</h2>
    {#if loadError}
      <p class="nav-error">{loadError}</p>
    {/if}
    {#if mailboxes != undefined}
      <div id="nav-mailboxes-list">
        {#each mailboxes as mailbox}
          <div
            class="nav-item"
            class:active={activeID === mailbox.id}
            role="button"
            tabindex="0"
            onclick={() => select(mailbox.id)}
            onkeydown={(e) => {
              if (e.key === "Enter" || e.key === " ") select(mailbox.id);
            }}
          >
            <span class="nav-item-name">{mailbox.username}</span>
            <button
              class="nav-item-sync"
              title="Sync mailbox"
              disabled={syncingID === mailbox.id}
              onclick={(e) => handleSync(e, mailbox.id)}
            >
              {syncingID === mailbox.id ? "…" : "↻"}
            </button>
          </div>
        {/each}
      </div>
    {:else if !loadError}
      <p class="nav-loading">Loading mailboxes…</p>
    {/if}
  </div>
  <a role="button" href="/mailboxes" id="manage-mailboxes-button">Manage Mailboxes</a>
</nav>
<main id="page-content">
  {@render children?.()}
</main>

<style>
  .nav-heading {
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: #767676;
    margin: 1rem 0 0.5rem;
  }

  .nav-error {
    color: #c0392b;
    font-size: 0.85rem;
    margin: 0.5rem 0;
  }

  .nav-loading {
    color: #999;
    font-size: 0.85rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
    cursor: pointer;
  }

  .nav-item-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .nav-item-sync {
    flex-shrink: 0;
    background: transparent;
    border: 1px solid rgba(0, 0, 0, 0.12);
    border-radius: 6px;
    color: var(--primary-color);
    font-size: 0.9rem;
    line-height: 1;
    padding: 2px 6px;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.15s ease;
  }

  .nav-item:hover .nav-item-sync,
  .nav-item.active .nav-item-sync {
    opacity: 1;
  }

  .nav-item-sync:disabled {
    cursor: wait;
    opacity: 1;
  }

  #page-content {
    flex: 1;
    height: 100vh;
    overflow: auto;
    display: flex;
    flex-direction: column;
  }
</style>
