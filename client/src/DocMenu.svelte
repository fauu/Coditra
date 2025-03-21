<script>
  let { onCloseClick } = $props();
</script>

<div class="doc-menu">
  <button aria-label="Close Document" onclick={onCloseClick}>
    <i class="la la-close"></i>
  </button>
</div>

<style>
  .doc-menu {
    position: fixed;
    top: 1rem;
    right: 1rem;
  }

  button {
    opacity: 0.1;
    font-size: 1.5rem;
    color: var(--color-fg);
    background: transparent;
    border: 1px solid var(--color-fg-lighter);
    border-radius: 4px;
    padding: 0.2rem 0.4rem;
    margin: 0;
  }

  button:hover,
  button:active {
    background: linear-gradient(
      30deg,
      var(--color-bg-darker),
      var(--color-bg) 40%,
      var(--color-bg-lighter-2) 100%
    );
    cursor: pointer;
  }

  button:hover {
    opacity: 1;
    border-color: var(--color-bg-darker-2);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    transform: translateY(-2px);
  }

  button:active {
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
    transform: translateY(0px);
  }
</style>
