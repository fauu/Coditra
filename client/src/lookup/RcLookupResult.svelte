<script>
  import TwoColumnTranslation from "./TwoColumnTranslation.svelte";

  let { lookupResult } = $props();
</script>

{#if lookupResult.translations}
  <div class="translations">
    {#each lookupResult.translations as translation}
      <!-- prettier-ignore -->
      <a class="translation further-lookup-button" href={translation.url} target="_blank">{translation.content}</a>
    {/each}
  </div>
{/if}

<div class="examples">
  {#each lookupResult.examples as example}
    <div class="example">
      <!-- NOTE: Removed 'rc' class here. Not sure if it was necessary. -->
      <TwoColumnTranslation object={example} />
    </div>
  {/each}
</div>

<style>
  .translations {
    margin: 0.75rem 0 0.85rem 0;
    display: flex;
    flex-flow: row wrap;
  }

  .translation {
    color: inherit;
    text-decoration: none;
    margin-bottom: 0.3rem;
  }

  .translation:not(:last-child) {
    margin-right: 0.3rem;
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

  :global(.source-rc em) {
    font-style: normal;
    background: linear-gradient(
      330deg,
      var(--color-hl) 30%,
      var(--color-hl-lighter) 100%
    );
  }
</style>
