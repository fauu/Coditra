<script>
  import WrEntryGroup from "./WrEntryGroup.svelte";

  export let lookupResult;
</script>

<div class="primary">
  {#if lookupResult.mainResults}
    <div class="group">
      <WrEntryGroup entries={lookupResult.mainResults} />
    </div>
  {/if}

  {#if lookupResult.compoundForms}
    <div class="group">
      <WrEntryGroup entries={lookupResult.compoundForms} />
    </div>
  {/if}
</div>

<div class="secondary">
  {#if lookupResult.otherDicts}
    {#each lookupResult.otherDicts as entry}
      <div class="other-dicts-entry">
        {@html entry}
      </div>
    {/each}
  {/if}
</div>

<style>
  .group:not(:last-child) {
    margin-bottom: 3rem;
  }

  .group:last-child {
    margin-bottom: 2rem;
  }

  :global(.other-dicts-entry ul) {
    list-style-type: none;
    margin-top: 0;
    padding-left: 2rem;
  }

  :global(.other-dicts-entry .headnumber) {
    font-weight: bold;
  }

  :global(.other-dicts-entry .subjarea) {
    font-variant: small-caps;
    text-transform: lowercase;
  }

  :global(.other-dicts-entry .ital) {
    font-style: italic;
  }

  :global(.other-dicts-entry .example.phrase) {
    border-bottom: 3px solid var(--color-bg-darker-2);
  }
</style>
