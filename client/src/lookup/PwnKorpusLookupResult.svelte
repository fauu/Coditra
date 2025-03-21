<script>
  let { lookupResult } = $props();
</script>

<div class="main">
  {#each lookupResult.entries as entry}
    <!-- prettier-ignore -->
    <div class="entry">
      {entry.preFragment}<a href={entry.detailsUrl} target="_blank">{entry.theWord}</a>{entry.postFragment}
    </div>
  {/each}
</div>

<style>
  .main:not(:last-child) {
    margin-bottom: 1.5rem;
  }

  .entry {
    border-bottom: 1px dashed var(--color-bg-darker-2);
    padding: 0.5rem 0;
  }

  .entry a {
    background: linear-gradient(
      330deg,
      var(--color-hl) 30%,
      var(--color-hl-lighter) 100%
    );
  }
</style>
