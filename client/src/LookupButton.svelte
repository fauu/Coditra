<script>
  import { createBubbler } from "svelte/legacy";

  const bubble = createBubbler();
  let { entry, isCurrent } = $props();
</script>

<button
  type="button"
  class="button"
  class:current={isCurrent}
  onclick={bubble("click")}
>
  <div class="content">
    <div class="sub-left">
      {#if entry.params && entry.params.sourceLang}
        <img
          class="flag"
          src="flags/{entry.params.sourceLang}.png"
          alt={entry.params.sourceLang}
        />
      {/if}
    </div>
    <!-- prettier-ignore -->
    <div class="main">
      <span class="name">{entry.name}</span>
      {#if entry.url}
        <i class="external-lookup-icon la la-external-link-alt"></i>
      {/if}
    </div>
    <div class="sub-right">
      {#if entry.params && entry.params.targetLang}
        <img
          class="flag"
          src="flags/{entry.params.targetLang}.png"
          alt={entry.params.targetLang}
        />
      {/if}
    </div>
  </div>
</button>

<style>
  .button {
    padding: 0.3rem 0.3rem;
    display: inline-block;
    background: var(--color-bg);
    background: linear-gradient(
      330deg,
      var(--color-bg-darker),
      var(--color-bg) 50%,
      var(--color-bg-lighter) 100%
    );
    cursor: pointer;
    border-right: 1px solid var(--color-bg-darker);
  }

  .button.current {
    pointer-events: none;
  }

  .button > .content {
    display: flex;
  }

  .button:hover,
  .button:active {
    background: linear-gradient(
      330deg,
      var(--color-bg-darker),
      var(--color-bg) 40%,
      var(--color-bg-lighter-3) 100%
    );
  }

  .button:active {
    box-shadow: inset 3px 5px 15px var(--color-bg-darker-3);
  }

  .button:active > .content {
    transform: translateY(1px);
  }

  .button {
    border-bottom: 1px solid var(--color-bg-darker);
  }

  .button.current {
    background: var(--color-bg-lighter);
    border-right: none;
  }

  .main {
    font-size: 0.8rem;
    font-weight: bold;
    flex: 1;
    line-height: 1.4rem;
  }

  .external-lookup-icon {
    font-size: 0.85rem;
    color: var(--color-fg-lighter);
  }

  .sub-left {
    width: 18px;
    margin-right: 0.35rem;
  }

  .flag {
    border: 1px solid var(--color-bg-darker);
  }
</style>
