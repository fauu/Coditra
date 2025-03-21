<script>
  let { object } = $props();
</script>

<div class="source">
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html object.source}
</div>
<div class="target">
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html object.target}
</div>

<style>
  .source,
  .target {
    flex: 1 0 0;
  }

  .source {
    padding-right: 1rem;
  }

  .target {
    padding-left: 1rem;
  }
</style>
