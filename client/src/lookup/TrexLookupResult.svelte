<script>
  import TwoColumnTranslation from "./TwoColumnTranslation.svelte";

  export let lookupResult;
</script>

{#if lookupResult.translations}
  <div class="translations">
    {#each lookupResult.translations as translation}
      <!-- prettier-ignore -->
      <a class="translation further-lookup-button" href="{translation.url}" target="_blank">{translation.content}</a>
    {/each}
  </div>
{/if}

{#if lookupResult.phrases}
  <div class="phrases">
    {#each lookupResult.phrases as phrase}
      <!-- prettier-ignore -->
      <a class="phrase further-lookup-button" href="{phrase.url}" target="_blank">{phrase.content}</a>
    {/each}
  </div>
{/if}

<div class="examples">
  {#each lookupResult.examples as example}
    <div class="example">
      <TwoColumnTranslation object={example} />
    </div>
  {/each}
</div>

<style>
  /* TIDY: Most of the styles repeat with RC */

  .translations,
  .phrases {
    margin: 0.75rem 0 0.85rem 0;
    display: flex;
    flex-flow: row wrap;
  }

  /* TIDY: Repeats with RC */
  .further-lookup-button {
    color: inherit;
    text-decoration: none;
    margin: 0 0.3rem 0.3rem 0;
  }

  .example {
    display: flex;
    margin-bottom: 1rem;
  }

  .example:first-child {
    padding-top: 0.7rem;
  }

  .example:not(:last-child) {
    border-bottom: 1px dashed var(--color-bg-darker-2);
    padding-bottom: 1rem;
  }

  :global(.source-trex em) {
    font-style: normal;
    background: linear-gradient(
      330deg,
      var(--color-hl) 30%,
      var(--color-hl-lighter) 100%
    );
  }
</style>
