<script>
  import { onMount } from "svelte";

  import { SENTENCE_SCROLL_TOP_OFFSET } from "./consts";
  import { scrollTo } from "./dom";

  export let doc;
  export let onMouseDown;
  export let extraMarginBottom;

  let segmentEls = [];

  onMount(() => hookSentences());

  const hookSentences = () => {
    if (!doc || segmentEls.length > 0) return;
    segmentEls = Array.from(document.getElementsByClassName("segment"));
    if (segmentEls.length === 0) {
      setTimeout(hookSentences, 100);
    } else {
      segmentEls.forEach((e) => {
        e.addEventListener("click", (event) => handleSentenceClick(event));
      });
    }
  };

  const handleSentenceClick = (event) => {
    const el = event.currentTarget;
    if (el.classList.contains("focused")) return;
    segmentEls.forEach((e) => (e.classList = "segment"));
    el.classList = "segment focused";
    scrollTo(el, { topOffset: SENTENCE_SCROLL_TOP_OFFSET, smooth: true });
  };
</script>

<div class="doc" on:mousedown={onMouseDown} style:margin-bottom={9 + extraMarginBottom}px>
  {@html doc.content}
</div>

<style>
  :global(.doc p) {
    border-bottom: 2px solid var(--color-bg-darker-3);
    margin: 1rem 0;
    padding-bottom: 1rem;
    cursor: default;
  }

  :global(.doc .segment) {
    display: block;
    padding: 0 0.3rem;
  }

  :global(.doc .segment:not(:last-child)) {
    margin-bottom: 1rem;
  }

  :global(.doc .segment.focused) {
    background: var(--color-bg-lighter);
    border-radius: 3px;
  }

  :global(.doc em) {
    border-bottom: 2px dotted var(--color-bg-darker-3);
  }

  .doc {
    margin: 9px 13px;
    font-size: 17px;
    line-height: 1.8;
    font-family: Merriweather, Georgia, serif;
    max-width: 960px; /* LookupPanel .side-column + .lookup-result */
  }
</style>
