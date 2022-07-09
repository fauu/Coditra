<script>
  import { SERVER_ADDR } from "../consts";

  export let large = false;
</script>

<div class="error" class:large>
  <!-- prettier-ignore -->
  <span>
    Could not connect to the server at <code>{SERVER_ADDR}</code>. Is it running?
  </span>
</div>

<style>
  .error {
    width: 100%;
    text-align: center;
  }

  .error.large {
    margin-top: 30vh;
    font-size: 1.1rem;
  }

  code {
    font-weight: bold;
  }

  .error.large code {
    font-size: 1.05rem;
  }
</style>
