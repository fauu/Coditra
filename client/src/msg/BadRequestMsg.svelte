<script>
  let { source } = $props();
</script>

<div class="error">
  <!-- prettier-ignore -->
  <span>
    Bad request. Is source <code>{source}</code> configured and implemented correctly?
  </span>
</div>

<style>
  .error {
    width: 100%;
    text-align: center;
  }

  code {
    font-size: 0.9rem;
    font-weight: bold;
  }
</style>
