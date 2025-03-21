<script>
  import TwoColumnTranslation from "./TwoColumnTranslation.svelte";

  let { object } = $props();
</script>

{#if object.expressions}
  <div class="expressions">
    {#each object.expressions as expression, i (i)}
      <div class="expression">
        <TwoColumnTranslation object={expression} />
      </div>
    {/each}
  </div>
{/if}

{#if object.examples}
  <div class="examples">
    {#each object.examples as example, i (i)}
      <div class="example">
        <TwoColumnTranslation object={example} />
      </div>
    {/each}
  </div>
{/if}

<style>
  .expression,
  .example {
    display: flex;
  }

  .expression:not(:last-child),
  .example:not(:last-child) {
    margin-bottom: 0.25rem;
    padding-bottom: 0.25rem;
    border-bottom: 1px dashed var(--color-bg-darker-2);
  }

  .expressions,
  .examples {
    margin: 1rem 0 2rem 0;
  }
</style>
