<script lang="ts">
  import type { Mail } from "$lib/types";
  import { getContext } from "svelte";
  const activeMailbox = getContext<{ id: number | undefined }>("activeMailbox");

  let mail = $state<Array<Mail> | undefined>(undefined);
  $effect(() => {
    if (activeMailbox.id == undefined) {
      return;
    }
    fetch(`/api/mailboxes/${activeMailbox.id}/messages`)
      .then((response) => response.json())
      .then((data) => {
        mail = data;
      });
  });
</script>

{#if activeMailbox.id === undefined}
  <div class="m-auto">
    <p>No mailbox selected.</p>
  </div>
{:else}
  <div style="flex flex-row">
    <div style="flex flex-col">
      <!-- list mail here -->
      {#if mail}
        {#each mail as message}
          {@render mailpreview(message)}
        {/each}
      {:else}
        <p>Loading mail.</p>
      {/if}
    </div>
  </div>
{/if}

{#snippet mailpreview(message)}
  <div>
    <h3>{message.subject}</h3>
  </div>
{/snippet}
