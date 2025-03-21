<script>
  import TwoColumnTranslation from "./TwoColumnTranslation.svelte";
  import LingeaExpressionsAndExamples from "./LingeaExpressionsAndExamples.svelte";

  let { lookupResult } = $props();
</script>

<div class="title-row">
  <span class="title">{lookupResult.title}</span>{#if lookupResult.femSuffix},
    <span class="fem-suffix">{lookupResult.femSuffix}</span>
  {/if}
</div>

<div class="morf-blocks">
  {#each lookupResult.morfBlocks as block}
    <div class="morf-block">
      <div class="morf">{block.morf}</div>
      <div class="morf-block-main">
        {#each block.segments as segment}
          <div class="segment">
            {#if segment.definition}
              <div class="definition">
                {@html segment.definition}
              </div>
            {/if}
            <LingeaExpressionsAndExamples object={segment} />
          </div>
        {/each}
      </div>
    </div>
  {/each}
</div>

{#if lookupResult.phrases}
  <LingeaExpressionsAndExamples object={lookupResult.phrases} />
{/if}

{#if lookupResult.keywordTerms}
  {#each lookupResult.keywordTerms as keywordTerm}
    <div class="keyword-term">
      <div class="keyword-term-main">
        <TwoColumnTranslation object={keywordTerm.twoColumn} />
      </div>
    </div>
  {/each}
{/if}

<!-- TODO: "related" links !-->

<style>
  .title-row {
    margin-bottom: 1rem;
  }

  .morf-blocks:not(:last-child) {
    margin-bottom: 1.25rem;
    padding-bottom: 0.5rem;
  }

  .morf-block {
    margin-bottom: 0.75rem;
  }

  .title,
  .fem-suffix {
    font-weight: bold;
    font-size: 1.1rem;
  }

  .morf {
    color: var(--color-fg-lighter);
    font-weight: bold;
  }

  .definition::before {
    content: "- ";
  }

  .keyword-term-main {
    display: flex;
  }

  .keyword-term:not(:last-child) {
    margin-bottom: 0.25rem;
    padding-bottom: 0.25rem;
    border-bottom: 1px dashed var(--color-bg-darker-2);
  }

  :global(.lex_ful_v, .lex_ful_w, .lex_ftx_w, .lex_ful_g) {
    margin: 0;
    color: var(--color-fg-lighter);
  }

  :global(.keyword-term-main .fh) {
    text-decoration: underline;
    margin: 0;
  }

  :global(.keyword-term-main .lex_ful_g) {
    font-size: 0.7rem;
  }

  :global(.keyword-term-main .lex_ful_g::before) {
    content: "(";
  }

  :global(.keyword-term-main .lex_ful_g::after) {
    content: ")";
  }

  :global(.keyword-term-main .lex_ftx_d) {
    font-style: italic;
    font-size: 0.9rem;
  }
</style>
