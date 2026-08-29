<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { api } from "$lib/api";
  import { fmtDateTime } from "$lib/mailutils";
  import type { Mail, Mailbox, Part } from "$lib/types";

  let mailboxId = $derived(Number(page.params.id));
  let mailId = $derived(page.params.message_id as string);

  let mailbox = $state<Mailbox | null>(null);
  let mail = $state<Mail | null>(null);
  let loading = $state(true);
  let error = $state("");

  async function load() {
    loading = true;
    error = "";
    try {
      const boxes = await api.listMailboxes();
      mailbox = boxes.find((b) => b.id === mailboxId) ?? null;
      mail = await api.getMail(mailboxId, mailId);
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load message";
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<div class="message-wrap">
  <div class="crumbs">
    <a class="crumb" href="/">box</a>
    <span class="sep">/</span>
    <a class="crumb" href={`/mailbox/${mailboxId}`}>
      {mailbox?.username ?? `Mailbox ${mailboxId}`}
    </a>
  </div>

  {#if loading}
    <p class="dim pad">Loading message…</p>
  {:else if error}
    <div class="error-banner">{error}</div>
  {:else if mail}
    <article class="card">
      <div class="head">
        <div class="avatar">{(mail.display_name || mail.from_address).slice(0, 1).toUpperCase()}</div>
        <div class="head-text">
          <span class="from">{mail.display_name || mail.from_address}</span>
          <span class="subject">{mail.subject || "(no subject)"}</span>
          <span class="date">
            {fmtDateTime(mail.recieved_at)}
            {#if mail.mailbox}
              <span class="inbox-tag">→ {mail.mailbox.username}</span>
            {/if}
          </span>
        </div>
      </div>
      <hr class="rule" />
      <div class="body-parts">
        {#each mail.body_parts as part (part.content_type)}
          {@render PartRenderer(part)}
        {/each}
      </div>
    </article>
    <div class="back">
      <a class="btn-ghost" href={`/mailbox/${mailboxId}`}>← Back to inbox</a>
    </div>
  {/if}
</div>

{#snippet PartRenderer(part: Part)}
  {#if part.content_type === "text/plain"}
    <pre class="body-text">{part.data ?? ""}</pre>
  {:else if part.content_type === "text/html"}
    <iframe class="body-html" title="Message content" srcdoc={part.data ?? ""} sandbox=""></iframe>
  {:else if part.content_type.startsWith("image/")}
    {#if part.link}
      <img class="body-image" src={part.link} alt="" role="presentation" />
    {:else if part.data}
      <img class="body-image" src={part.data} alt="" role="presentation" />
    {/if}
  {:else if part.link}
    <a class="attachment" href={part.link} download>
      📎 {part.content_type}
    </a>
  {:else}
    <div class="attachment-meta">Attachment: {part.content_type}</div>
  {/if}
{/snippet}

<style>
  .message-wrap {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 720px;
    margin: 0 auto;
    padding: 1.5rem;
  }

  .crumbs {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-weight: 600;
    font-size: 0.9rem;
    color: #767676;
    margin-bottom: 1rem;
  }

  .crumb {
    color: var(--primary-color);
    text-decoration: none;
  }

  .sep {
    opacity: 0.5;
  }

  .dim {
    color: #767676;
  }

  .pad {
    padding: 1rem 0;
  }

  .error-banner {
    background: #fdecea;
    color: #a33;
    border: 1px solid #f5c6c0;
    border-radius: 8px;
    padding: 0.6rem 0.9rem;
    font-size: 0.9rem;
  }

  .card {
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 12px;
    padding: 1.5rem;
    background: #fff;
  }

  .head {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .avatar {
    width: 48px;
    height: 48px;
    flex-shrink: 0;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--primary-color), #8a5bef);
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
    font-size: 0.95rem;
    font-weight: 650;
  }

  .subject {
    font-size: 1.25rem;
    font-weight: 700;
    margin: 0.1rem 0;
  }

  .date {
    color: #959595;
    font-size: 0.8rem;
  }

  .inbox-tag {
    margin-left: 0.5rem;
    color: var(--primary-color-dark);
  }

  .rule {
    border: none;
    border-top: 1px solid rgba(0, 0, 0, 0.08);
    margin: 1.25rem 0;
  }

  .body-parts {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .body-text {
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.6;
    font-size: 0.95rem;
    font-family: inherit;
    margin: 0;
  }

  .body-html {
    width: 100%;
    border: 1px solid rgba(0, 0, 0, 0.08);
    border-radius: 8px;
    min-height: 200px;
  }

  .body-image {
    max-width: 100%;
    border-radius: 8px;
  }

  .attachment {
    display: inline-block;
    padding: 0.5rem 0.75rem;
    border: 1px solid rgba(0, 0, 0, 0.12);
    border-radius: 8px;
    color: var(--primary-color-dark);
    text-decoration: none;
  }

  .attachment-meta {
    color: #959595;
    font-size: 0.85rem;
  }

  .back {
    margin-top: 1.25rem;
  }

  .btn-ghost {
    color: var(--primary-color-dark);
    text-decoration: none;
    font-size: 0.9rem;
  }

  .btn-ghost:hover {
    text-decoration: underline;
  }
</style>
