<script>
  export let object;
</script>

<div class="source">
  {@html object.source}
</div>
<div class="target">
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
