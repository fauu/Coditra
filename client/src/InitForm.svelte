<script>
  export let docNames;
  export let onSubmit;
</script>

<div class="main">
  <h1>Coditra</h1>
  <p class="select-doc-prompt-msg">Select a document:</p>
  <ul class="doc-list">
    {#each docNames as docName}
      <li on:click={() => onSubmit(docName)}>{docName}</li>
    {/each}
  </ul>
  <br>
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

  .doc-list li {
    padding: 0.4rem 0.55rem;
    background-color: var(--color-bg-lighter);
  }

  .doc-list li:first-child {
    border-radius: 4px 4px 0 0; 
  }

  .doc-list li:last-child {
    border-radius: 0 0 4px 4px; 
  }

  .doc-list li:not(:last-child) {
    border-bottom: 1px solid var(--color-bg-darker-2);
  }

  .doc-list li:hover {
    background-color: var(--color-bg-lighter-2);
    cursor: pointer;
  }
</style>
