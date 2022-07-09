<script>
  export let entries;
</script>

{#each entries as entry}
  <div class="entry">
    <div class="source-column">
      <div class="source-form">{entry.sourceForm}</div>
      {#if entry.sense}
        <div class="sense">
          {#if entry.sense.note}
            <span class="sense-note">{entry.sense.note}</span>
          {/if}
          ({entry.sense.sense})
        </div>
      {/if}
    </div>
    <div class="target-column">
      {#each entry.translations as translation}
        <div class="translation-row">
          <div class="note">
            {#if translation.note}({translation.note}){/if}
          </div>
          <div class="translation">{translation.translation}</div>
        </div>
      {/each}
    </div>
  </div>
{/each}

<style>
  .entry {
    display: flex;
  }

  .entry:first-child {
    margin-top: 0.5rem;
  }

  .entry:not(:last-child) {
    margin-bottom: 0.75rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px dashed var(--color-bg-darker-2);
  }

  .source-column,
  .target-column {
    display: flex;
    flex-direction: column;
  }

  .source-column {
    flex: 4;
  }

  .target-column {
    flex: 6;
  }

  .translation-row {
    display: flex;
    flex-direction: row;
  }

  .source-form {
    font-weight: bold;
  }

  .sense {
    font-size: 0.85rem;
  }

  .sense-note {
    font-style: italic;
  }

  .note {
    flex: 1;
    text-align: right;
    font-style: italic;
    padding-right: 0.35rem;
    font-size: 0.85rem;
    color: var(--color-fg-lighter);
  }

  .translation {
    flex: 2;
  }
</style>
