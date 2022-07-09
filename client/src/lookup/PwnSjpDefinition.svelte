<script>
  export let definition;
  export let onRefLookup;
</script>

{#if definition.qualifier}
  <span class="qualifier">{definition.qualifier}</span>
{/if}
{#if definition.xref}
  <span>zob.</span>
  <button class="further-lookup-button" on:click={onRefLookup(definition.xref)}>
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
