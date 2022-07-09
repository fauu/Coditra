<script>
  export let delayMs = 1000;

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
  <div class="rect1" />
  <div class="rect2" />
  <div class="rect3" />
  <div class="rect4" />
  <div class="rect5" />
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
