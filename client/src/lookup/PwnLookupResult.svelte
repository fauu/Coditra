<script>
  import PwnSjpDefinition from "./PwnSjpDefinition.svelte";

  export let lookupResult;
  export let onRefLookup;
</script>

<div class="sjp-result">
  {#each lookupResult.sjpResult.entries || [] as entry}
    <div class="entry">
      <span class="entry-title">{entry.title}</span>
      {#if entry.definitions.length === 1}
        <PwnSjpDefinition definition={entry.definitions[0]} {onRefLookup} />
      {:else}
        <ol class="definition-list">
          {#each entry.definitions as definition}
            <li>
              <PwnSjpDefinition {definition} {onRefLookup} />
            </li>
          {/each}
        </ol>
      {/if}
    </div>
  {/each}
</div>

<div class="doroszewski-result">
  {#each lookupResult.doroszewskiResult.entries || [] as entry}
    <div class="entry">
      {#each entry.imgFragmentUrls as imgFragmentUrl}
        <img
          class="doroszewski-image-fragment"
          src={imgFragmentUrl}
          alt={entry.title}
        />
      {/each}
    </div>
  {/each}
</div>

<style>
  .sjp-result:not(:last-child) {
    margin-bottom: 1.5rem;
  }

  .sjp-result .entry {
    margin-bottom: 0.75rem;
  }

  .doroszewski-result {
    display: flex;
    flex-direction: column;
  }

  .doroszewski-result .entry {
    margin-bottom: 0.5rem;
  }

  .doroszewski-result .entry {
    align-self: flex-start;
    border: 1px solid var(--color-bg-darker);
    font-size: 0; /* Gets rid of unwanted space between Doroszewski fragments */
    background: #ffffff;
    padding: 0.5rem;
  }

  .doroszewski-image-fragment {
    display: block;
  }

  .entry-title {
    font-weight: bold;
    margin-right: 0.15rem;
  }

  .definition-list {
    margin: 0;
    padding-left: 1.8rem;
    line-height: 1.2;
  }
</style>
