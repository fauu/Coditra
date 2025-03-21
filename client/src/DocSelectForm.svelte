<script>
  let { docNames, onSubmit } = $props();
</script>

<div class="main">
  <h1>Coditra</h1>
  <p class="select-doc-prompt-msg">Select a document:</p>
  <ul class="doc-list">
    {#each docNames as docName}
      <button type="button" onclick={() => onSubmit(docName)}>{docName}</button>
    {/each}
  </ul>
  <br />
</div>

<style>
  .main {
    margin: 3rem auto 1rem;
  }

  h1 {
    color: var(--color-fg-lighter);
    letter-spacing: -0.05rem;
  }

  .select-doc-prompt-msg {
    font-size: 1.1rem;
    color: var(--color-fg-lighter);
  }

  .doc-list {
    list-style-type: none;
    margin: 0;
    padding: 0;
    font-family: var(--fonts-mono);
    border: 1px solid var(--color-bg-darker-2);
    border-radius: 4px;
  }

  .doc-list button {
    padding: 0.4rem 0.55rem;
    background-color: var(--color-bg-lighter);
    width: 100%;
  }

  .doc-list button:first-child {
    border-radius: 4px 4px 0 0;
  }

  .doc-list button:last-child {
    border-radius: 0 0 4px 4px;
  }

  .doc-list button:not(:last-child) {
    border-bottom: 1px solid var(--color-bg-darker-2);
  }

  .doc-list button:hover {
    background-color: var(--color-bg-lighter-2);
    cursor: pointer;
  }
</style>
