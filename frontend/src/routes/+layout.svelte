<script lang="ts">
  import "./layout.css";
  import type { Mailbox } from "$lib/types";
  import favicon from "$lib/assets/favicon.svg";
  import { setContext } from "svelte";

  let { children } = $props();
  let mailboxes = $state<Mailbox[] | undefined>(undefined);
  let activeMailboxID = $state<number | undefined>(undefined);

  async function getMailboxes() {
    const resp = await fetch("/api/mailboxes");
    return await resp.json();
  }
  getMailboxes().then((data) => {
    mailboxes = data;
  });

  setContext("activeMailbox", {
    get id() {
      return activeMailboxID;
    },
  });
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<nav>
  <div>
    <h1 class="text-center">box</h1>
    <h2>Mailboxes</h2>
    {#if mailboxes != undefined}
      <div id="nav-mailboxes-list">
        {#each mailboxes as mailbox}
          <button
            class="nav-mailbox text-center"
            onclick={() => {
              activeMailboxID = mailbox.id;
            }}
            class:active={activeMailboxID === mailbox.id}>{mailbox.username}</button
          >
        {/each}
        <button
          class="nav-mailbox"
          onclick={() => {
            activeMailboxID = 2;
          }}
          class:active={activeMailboxID === 2}>meow</button
        >
      </div>
    {/if}
  </div>
  <a role="button" href="/mailboxes" id="manage-mailboxes-button">Manage Mailboxes</a>
</nav>
{@render children({ activeMailboxID })}
