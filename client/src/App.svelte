<script>
  import * as api from "./api.js";

  import { 
    LOOKUP_PANEL_MIN_HEIGHT,
    LOOKUP_PANEL_MIN_TOP_MARGIN,
    LOOKUP_PANEL_DEFAULT_HEIGHT
  } from "./consts";

  import Doc from "./Doc.svelte";
  import DocMenu from "./DocMenu.svelte";
  import InitForm from "./InitForm.svelte";
  import LookupPanel from "./LookupPanel.svelte";
  import Spinner from "./Spinner.svelte";

  let appState = "init-setup";
  let showLookupPanel = false;
  let selectedDocName = localStorage.getItem("docName") || undefined;
  let lookupPanelHeight = LOOKUP_PANEL_DEFAULT_HEIGHT;
  let isResizingLookupPanel = false;
  let doc;
  let setups = api.getSetups();
  let mouseUpSelectionTimeout;

  const savedLookupPanelHeight = parseInt(localStorage.getItem("lookupPanelHeight"));
  if (!isNaN(savedLookupPanelHeight)) {
    lookupPanelHeight = savedLookupPanelHeight;
  }

  const viewSelectedDoc = async () => {
    appState = "preparing-doc-view";
    try {
      const [d, s] = await Promise.all([api.getDoc(selectedDocName), setups]);
      doc = d;
      setups = s;
      appState = "doc-view";
    } catch (_err) {
      localStorage.removeItem("docName");
      appState = "init-setup";
    }
  };

  if (selectedDocName !== undefined) {
    viewSelectedDoc();
  }

  const closeDoc = () => {
    selectedDocName = undefined;
    doc = undefined;
    localStorage.removeItem("docName");
    appState = "init-setup";
    showLookupPanel = false;
  };

  const handleInitSetupSubmit = async (docName) => {
    selectedDocName = docName;
    localStorage.setItem("docName", selectedDocName);
    viewSelectedDoc();
  };

  const handleCloseDocClick = () => {
    closeDoc();
  };

  const handleWindowMouseMove = (e) => {
    if (isResizingLookupPanel && e.movementY !== 0) {
      lookupPanelHeight -= e.movementY;
      const maxHeight = window.innerHeight - LOOKUP_PANEL_MIN_TOP_MARGIN;
      if (lookupPanelHeight < LOOKUP_PANEL_MIN_HEIGHT) {
        lookupPanelHeight = LOOKUP_PANEL_MIN_HEIGHT;
      } else if (lookupPanelHeight > maxHeight) {
        lookupPanelHeight = maxHeight;
      }
    }
  };

  const handleWindowMouseUp = () => {
    if (isResizingLookupPanel) {
      isResizingLookupPanel = false;
      document.body.style.cursor = "default";
      localStorage.setItem("lookupPanelHeight", lookupPanelHeight);
    }

    const selection = document.getSelection();
    const selectionText = selection.toString();
    if (selection.type === "Caret") return;
    mouseUpSelectionTimeout = setTimeout(() => {
      if (selection.type === "Caret" || selectionText.trim() === "") return;
      if (appState === "doc-view") {
        showLookupPanel = true;
      }
    }, 175);
  };

  const handleWindowKeyUp = (ev) => {
    switch (ev.key) {
      case "l":
        if (appState === "doc-view") showLookupPanel = true;
    }
  };

  const handleDocMouseDown = () => {
    clearTimeout(mouseUpSelectionTimeout);
    showLookupPanel = false;
  };
  
  const handleLookupPanelResizeStart = () => {
    document.body.style.cursor = "row-resize";
    isResizingLookupPanel = true;
  };
</script>

<svelte:window 
  on:mouseup={handleWindowMouseUp} 
  on:keyup={handleWindowKeyUp} 
  on:mousemove={handleWindowMouseMove} />

<main class="main-container">
  {#if appState === "init-setup"}
    {#await api.getDocNames()}
      <Spinner />
    {:then docNames}
      <InitForm
        {docNames}
        onSubmit={handleInitSetupSubmit}
      />
    {:catch error}
      {error}
    {/await}
  {:else if appState === "preparing-doc-view"}
    <Spinner />
  {:else if appState === "doc-view"}
    <DocMenu onCloseClick={handleCloseDocClick} />
    <Doc 
      {doc} 
      onMouseDown={handleDocMouseDown}
      bind:extraMarginBottom={lookupPanelHeight} />
    {#if showLookupPanel}
      <LookupPanel
        {setups}
        selectedText={document.getSelection().toString()}
        bind:height={lookupPanelHeight}
        onResizeStart={handleLookupPanelResizeStart}
      />
    {/if}
  {/if}
</main>

<style>
  :root {
    --color-fg: #1b1b1a;
    --color-fg-lighter: #644f3b;
    --color-fg-lighter-2: #9d7e5e;
    --color-bg: #e6d0b9;
    --color-bg-darker: #dbc0a4;
    --color-bg-darker-2: #cdaf90;
    --color-bg-darker-3: #bf9e7c;
    --color-bg-lighter: #ecdac8;
    --color-bg-lighter-2: #f6ebe1;
    --color-bg-lighter-3: #fffcf9;
    --color-hl: #ffbaef;
    --color-hl-darker: #eeaddf;
    --color-hl-darker-3: #b95ea4;
    --color-hl-darker-4: #651782;
    --color-hl-lighter: #ffc7f2;
    --fonts-default: -apple-system, BlinkMacSystemFont, Roboto, "Droid Sans", "Helvetica Now",
      "Helvetica Neue", Helvetica, Geneva, Arial, sans-serif;
    --fonts-mono: "Roboto Mono", "Consolas", "Droid Sans Mono", monospace;
  }

  :global(html) {
    box-sizing: border-box;
    font-size: 100%;
    font-family: var(--fonts-default);
  }

  :global(*),
  :global(*:before),
  :global(*:after) {
    box-sizing: inherit;
  }

  :global(body) {
    margin: 0;
    padding: 0;
    background: var(--color-bg);
    color: var(--color-fg);
  }

  :global(a),
  :global(a:active),
  :global(a:visited) {
    color: var(--color-hl-darker-4);
  }

  .main-container {
    width: 100%;
    display: flex;
  }
</style>
