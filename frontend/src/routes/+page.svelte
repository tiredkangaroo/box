<script lang="ts">
  import { goto } from "$app/navigation";
  import { api } from "$lib/api";
  import type { Mailbox } from "$lib/types";
  import { onMount } from "svelte";

  let mailboxes = $state<Mailbox[]>([]);
  let loaded = $state(false);

  onMount(async () => {
    try {
      mailboxes = await api.listMailboxes();
    } finally {
      loaded = true;
    }
  });
</script>

<div class="home-wrap">
  <h1>box</h1>
  {#if !loaded}
    <p class="dim">Loading mailboxes…</p>
  {:else if mailboxes.length === 0}
    <p class="dim">No mailboxes yet. Add one from the sidebar to get started.</p>
  {:else}
    <p class="dim">Pick a mailbox from the sidebar to view your mail, or manage them below.</p>
    <a class="btn" href="/mailboxes">Manage Mailboxes</a>
  {/if}
</div>

<style>
  .home-wrap {
    margin: auto;
    text-align: center;
    padding: 2rem;
  }

  .dim {
    color: #767676;
  }

  .btn {
    display: inline-block;
    margin-top: 1rem;
    background: linear-gradient(180deg, var(--primary-color), var(--primary-color-dark));
    border: 1px solid #308cb0;
    color: white;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
    padding: 0.5rem 1rem;
    border-radius: 10px;
    font-size: 0.9rem;
    cursor: pointer;
    text-decoration: none;
  }
</style>
