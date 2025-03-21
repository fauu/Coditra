<script>
  let { lookupResult, onRefLookup } = $props();
</script>

{#if lookupResult.title}
  <div class="header">
    Synonimy słowa
    <span class="title">{lookupResult.title}</span>
  </div>

  {#each lookupResult.synonymGroups as synonymGroup}
    <div class="synonym-group">
      <ul class="synonym-list">
        {#each synonymGroup as synonym}
          <li class="synonym">
            {#if synonym.extra}
              <span class="extra">({synonym.extra})</span>
            {/if}
            {synonym.synonym}
          </li>
        {/each}
      </ul>
    </div>
  {/each}
{:else}
  <div class="header">
    Nie znaleziono synonimów podanego słowa. Proponowane alternatywy:
  </div>

  <ul class="synonym-list">
    {#each lookupResult.suggestedAlternatives as alt}
      <li class="synonym">
        <button class="further-lookup-button" onclick={() => onRefLookup(alt)}>
          {alt}
        </button>
      </li>
    {/each}
  </ul>
{/if}

<style>
  .header {
    margin-bottom: 1rem;
    font-size: 1.1rem;
  }

  .title {
    font-weight: bold;
  }

  .synonym-group:not(:last-child) {
    margin-bottom: 1rem;
  }

  .extra {
    font-style: italic;
    font-size: 0.9rem;
  }

  .synonym-list {
    margin: 0;
    padding: 0;
    list-style-type: none;
  }

  .synonym-list > * + * {
    margin-top: 0.25rem;
  }

  .synonym-list li::before {
    content: "– ";
  }
</style>
