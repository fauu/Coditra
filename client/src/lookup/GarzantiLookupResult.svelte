<script>
  let { lookupResult } = $props();
</script>

<div class="title-row">
  <span class="title">
    <!-- eslint-disable-next-line svelte/no-at-html-tags -->
    {@html lookupResult.word}
  </span>
  <span class="hyphenation">{lookupResult.hyphenation}</span>
</div>

<div class="definition">
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html lookupResult.definitionHtml}
</div>

<style>
  .title-row {
    margin-bottom: 1rem;
  }

  .title {
    font-weight: bold;
    font-size: 1.1rem;
  }

  :global(.definition p) {
    margin: 0.5rem 0;
  }

  :global(.definition a) {
    color: var(--color-fg-lighter);
  }
</style>
