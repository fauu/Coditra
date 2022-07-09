<script>
  import { onMount } from "svelte";

  import BadRequestMsg from "./msg/BadRequestMsg.svelte";
  import NoResultsMsg from "./msg/NoResultsMsg.svelte";
  import ServerConnectionErrorMsg from "./msg/ServerConnectionErrorMsg.svelte";

  import * as api from "./api.js";
  import LookupButton from "./LookupButton.svelte";
  import Spinner from "./Spinner.svelte";

  import PwnLookupResult from "./lookup/PwnLookupResult.svelte";
  import PwnKorpusLookupResult from "./lookup/PwnKorpusLookupResult.svelte";
  import WrLookupResult from "./lookup/WrLookupResult.svelte";
  import RcLookupResult from "./lookup/RcLookupResult.svelte";
  import LingeaLookupResult from "./lookup/LingeaLookupResult.svelte";
  import SynonimyplLookupResult from "./lookup/SynonimyplLookupResult.svelte";
  import GarzantiLookupResult from "./lookup/GarzantiLookupResult.svelte";
  import TrexLookupResult from "./lookup/TrexLookupResult.svelte";

  export let setups;
  export let selectedText;
  export let height;
  export let onResizeStart;

  const sourceToLookupResultComponent = {
    pwn: PwnLookupResult,
    pwnkorpus: PwnKorpusLookupResult,
    wr: WrLookupResult,
    rc: RcLookupResult,
    lingea: LingeaLookupResult,
    synonimypl: SynonimyplLookupResult,
    garzanti: GarzantiLookupResult,
    trex: TrexLookupResult,
  };

  let selectedSetupIdx = 0;
  let selectedSetup;
  let inputEl;
  let input;
  let lookupResultComponent;
  let currentLookupEntry;
  let lookupFetch = { state: "initial" };
  let defaultLookupRunning;

  $: selectedSetup = setups[selectedSetupIdx];
  $: {
    if (selectedSetup) {
      localStorage.setItem("setupName", selectedSetup.name);
    }
  }
  const savedSetupName = localStorage.getItem("setupName");
  if (savedSetupName) {
    const foundIdx = setups.findIndex((s) => s.name === savedSetupName);
    if (foundIdx !== -1) {
      selectedSetupIdx = foundIdx;
    }
  }

  onMount(() => {
    input = selectedText;
    inputEl.focus();
    if (input !== "") {
      const firstDefaultLookupEntry = selectedSetup.lookupEntries.find(
        (e) => e.default
      );
      if (firstDefaultLookupEntry) {
        defaultLookupRunning = true;
        doLookup(firstDefaultLookupEntry, input);
      }
    }
  });

  $: input && resetLookup();

  const resetLookup = () => {
    if (lookupFetch.state !== "initial" && !defaultLookupRunning)
      lookupFetch = { state: "initial" };
  };

  const doLookup = async (entry, input) => {
    if (!entry.source) {
      window.open(entry.url.replace("{input}", input));
      return;
    }
    if (lookupFetch.state === "fetching") {
      // TODO: Should probably cancel the ongoing request too if possible
      resetLookup();
    }
    lookupResultComponent = sourceToLookupResultComponent[entry.source];
    currentLookupEntry = entry;
    lookupFetch.state = "fetching";
    try {
      lookupFetch.result = await api.getDef(entry.source, input, entry.params);
      // Make sure we haven't cancelled in the meantime
      if (lookupFetch.state === "fetching") {
        lookupFetch.state = "done";
      }
    } catch (err) {
      if (lookupFetch.state !== "fetching") {
        return;
      }
      lookupFetch.state = "error";
      if (err.message.startsWith("Request failed")) {
        lookupFetch.error = "lookup";
      } else {
        lookupFetch.error = "server-connection";
      }
    }
    defaultLookupRunning = false;
  };

  const handleMouseDown = (e) => {
    if (e.button === 0 && e.offsetY < 5) {
      onResizeStart();
    }
  };

  const handleLookupButtonClick = async (entry) => {
    if (!input) return;
    doLookup(entry, input);
  };
</script>

<div class="main" on:mousedown={handleMouseDown} style:height={height}px>
  <div class="input-row">
    <input class="input" type="text" bind:this={inputEl} bind:value={input} />
  </div>
  <div class="main-row">
    <div class="side-column">
      <div class="lookup-buttons">
        {#each selectedSetup.lookupEntries as entry}
          <LookupButton
            {entry}
            isCurrent={currentLookupEntry === entry &&
              lookupFetch.state === "done"}
            on:click={() => handleLookupButtonClick(entry)}
          />
        {/each}
      </div>
      <select class="setup-picker" bind:value={selectedSetupIdx} required>
        {#each setups as setup, idx}
          <option value={idx}>{setup.name}</option>
        {/each}
      </select>
    </div>
    <div
      class="lookup-result {currentLookupEntry
        ? 'source-' + currentLookupEntry.source
        : ''}"
      class:loaded={lookupFetch.state === "done"}
    >
      {#if lookupFetch.state === "fetching"}
        <div class="spinner-container">
          <Spinner delayMs={1000} />
        </div>
      {:else if lookupFetch.state === "done"}
        {#if !lookupFetch.result || lookupFetch.result.isEmpty}
          <div class="msg-container">
            <NoResultsMsg />
          </div>
        {:else}
          <div class="lookup-result-source-link">
            <a href={lookupFetch.result.sourceUrl} target="_blank">
              {lookupFetch.result.sourceUrl}
            </a>
          </div>
          <svelte:component
            this={lookupResultComponent}
            lookupResult={lookupFetch.result}
            onRefLookup={(input) => doLookup(currentLookupEntry, input)}
          />
        {/if}
      {:else if lookupFetch.state === "error"}
        <div class="msg-container">
          {#if lookupFetch.error === "lookup"}
            <BadRequestMsg source={currentLookupEntry.source} />
          {:else}
            <ServerConnectionErrorMsg />
          {/if}
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .main {
    position: fixed;
    bottom: 0;
    width: 100%;
    background: var(--color-bg);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
  }

  .main:before {
    content: "";
    position: absolute;
    height: 5px;
    width: 100%;
    cursor: row-resize;
  }

  .input-row {
    border-bottom: 1px solid var(--color-bg-darker);
  }

  .input {
    margin: 0;
    width: 100%;
    background: var(--color-bg-lighter-3);
    font-size: 0.9rem;
    font-family: var(--fonts-mono);
    font-weight: 500;
    border: none;
    border-radius: 4px 4px 0 0;
    padding: 0.45rem 0.5rem;
    height: 35px;
  }

  .input:focus {
    outline: none;
    box-shadow: inset 0px 0px 5px 0px var(--color-hl);
  }

  .main-row {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .side-column {
    flex: 0 0 160px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .lookup-buttons {
    display: flex;
    flex-direction: column;
  }

  .setup-picker {
    margin-bottom: 0;
    background: var(--color-bg-lighter-2);
    border: 1px solid var(--color-bg-darker);
    border-bottom: 0;
    border-left: 0;
    font-size: 0.85rem;
  }

  .lookup-result {
    padding: 0.35rem 0.8rem;
    flex: 1;
    overflow-y: auto;
    font-size: 0.9rem;
    position: relative;
    max-width: 826px; /* 960px - .side-column + 2 * .doc horizontal padding */
  }

  :global(.lookup-result > :last-child) {
    margin-bottom: 0.35rem;
  }

  .lookup-result.loaded {
    background: linear-gradient(
      90deg,
      var(--color-bg-lighter),
      var(--color-bg) 30%
    );
  }

  .lookup-result-source-link {
    font-size: 0.8rem;
    margin-bottom: 0.5rem;
  }

  .spinner-container, .msg-container {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0;
  }

  :global(.further-lookup-button) {
    all: unset;
    cursor: pointer;
    border: 1px dashed var(--color-fg-lighter-2);
    padding: 0.1rem 0.24rem;
    border-radius: 2px;
  }

  :global(.further-lookup-button:hover) {
    background-color: var(--color-hl);
  }

  :global(.further-lookup-button:focus) {
    outline: 1px dotted var(--color-fg-lighter);
    outline-offset: 2px;
  }

  :global(.further-lookup-button:active) {
    background-color: var(--color-hl-darker) !important;
  }
</style>
