<script>
  /**
   * @typedef {Object} Props
   * @property {number} [delayMs]
   */

  /** @type {Props} */
  let { delayMs = 1000 } = $props();

  const delayedReveal = (node) => {
    node.style.display = "none";
    setTimeout(() => {
      node.style.display = "initial";
    }, delayMs);
    return {};
  };
</script>

<!-- https://github.com/tobiasahlin/SpinKit/blob/master/LICENSE -->
<div use:delayedReveal class="spinner">
  <div class="rect1"></div>
  <div class="rect2"></div>
  <div class="rect3"></div>
  <div class="rect4"></div>
  <div class="rect5"></div>
</div>

<style>
  .spinner {
    margin: 100px auto;
    width: 50px;
    height: 40px;
    text-align: center;
    font-size: 10px;
  }

  .spinner > div {
    background-color: var(--color-fg-lighter-2);
    height: 100%;
    width: 6px;
    display: inline-block;

    animation: sk-stretchdelay 1.2s infinite ease-in-out;
  }

  .spinner .rect2 {
    animation-delay: -1.1s;
  }

  .spinner .rect3 {
    animation-delay: -1s;
  }

  .spinner .rect4 {
    animation-delay: -0.9s;
  }

  .spinner .rect5 {
    animation-delay: -0.8s;
  }

  @keyframes sk-stretchdelay {
    0%,
    40%,
    100% {
      transform: scaleY(0.4);
    }
    20% {
      transform: scaleY(1);
    }
  }
</style>
