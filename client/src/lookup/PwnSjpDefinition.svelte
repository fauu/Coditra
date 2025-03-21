<script>
  let { definition, onRefLookup } = $props();
</script>

{#if definition.qualifier}
  <span class="qualifier">{definition.qualifier}</span>
{/if}
{#if definition.xref}
  <span>zob.</span>
  <button
    class="further-lookup-button"
    onclick={() => onRefLookup(definition.xref)}
  >
    {definition.content}
  </button>
{:else}
  <span class="content">«{definition.content}»</span>
{/if}

<style>
  .qualifier {
    font-style: italic;
  }
</style>
