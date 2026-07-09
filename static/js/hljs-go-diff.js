// Registers a `go-diff` language on highlight.js that renders a unified diff
// with Go syntax highlighting inside every +/-/context line.
//
// Implementation note: we don't compose this out of hljs `contains` rules
// with `subLanguage: 'go'`. That approach - which is what the built-in
// `diff` language does structurally - misses the first `+` line after any
// context region when the sub-language is set, because hljs's contains
// loop leaves the parser mid-line when a subLanguage region ends and the
// `^`-anchored diff regex fails to re-anchor at the next line start.
//
// Instead, we highlight the whole block as Go first, then walk the
// resulting DOM in a `after:highlightElement` hook and wrap every physical
// line whose first character is `+` or `-` in an `hljs-addition` or
// `hljs-deletion` span. Depends on highlight.js already having the
// built-in `go` language loaded.
(function () {
  function register() {
    if (!window.hljs || typeof window.hljs.registerLanguage !== 'function') {
      // code.min.js hasn't finished evaluating yet; try again on the next tick.
      return setTimeout(register, 0);
    }

    // Register go-diff as an alias for go so hljs happily highlights the
    // whole block as Go. The diff-line wrapping happens after the fact.
    window.hljs.registerLanguage('go-diff', function (hljs) {
      var go = hljs.getLanguage('go');
      if (!go) {
        // Fall back to plaintext so we at least render something legible.
        return { name: 'Go Diff (go not loaded)', contains: [] };
      }
      // Return a shallow clone so we don't mutate the shared `go` mode.
      var mode = Object.assign({}, go);
      mode.name = 'Go Diff';
      mode.aliases = ['diff-go', 'godiff'];
      return mode;
    });

    // After hljs has highlighted a code block, wrap +/- lines with the
    // diff classes. This runs for every element hljs processes; we no-op
    // unless it's a `language-go-diff` block.
    window.hljs.addPlugin({
      'after:highlightElement': function (data) {
        var el = data.el;
        if (!el || !el.classList || !el.classList.contains('language-go-diff')) {
          return;
        }
        wrapDiffLines(el);
      },
    });

    function wrapDiffLines(codeEl) {
      // We walk the code element's child nodes in order, tracking the
      // current physical line's DOM range. When we hit a text node
      // containing a newline, we finalise the line (wrap it if its first
      // char was + or -) and start a new one.
      var lineNodes = [];
      var firstChar = null;

      function flushLine(trailingNewlineTextNode) {
        if (!lineNodes.length && firstChar == null) return;
        var cls = firstChar === '+' ? 'hljs-addition'
                : firstChar === '-' ? 'hljs-deletion'
                : null;
        if (cls) {
          var wrap = document.createElement('span');
          wrap.className = cls;
          // Insert wrapper before the first line node, then move all line
          // nodes into it.
          lineNodes[0].parentNode.insertBefore(wrap, lineNodes[0]);
          for (var i = 0; i < lineNodes.length; i++) {
            wrap.appendChild(lineNodes[i]);
          }
        }
        lineNodes = [];
        firstChar = null;
      }

      // Flatten the traversal: for each direct or nested descendant, we
      // want to attribute it to a line. But wrapping across arbitrary DOM
      // nesting is fragile - instead, we work at the direct-child level
      // and split any child text node that contains newlines.
      splitTextNodesOnNewlines(codeEl);

      var child = codeEl.firstChild;
      while (child) {
        var next = child.nextSibling;
        if (child.nodeType === Node.TEXT_NODE && child.nodeValue === '\n') {
          // End of line - flush before consuming the newline.
          flushLine();
          child = next;
          continue;
        }
        // Record the first non-newline character of this line if not yet set.
        if (firstChar == null) {
          var t = child.textContent || '';
          if (t.length > 0) {
            firstChar = t.charAt(0);
          }
        }
        lineNodes.push(child);
        child = next;
      }
      // Trailing line (no terminating newline).
      flushLine();
    }

    // Split every direct-child text node of `parent` at newline boundaries
    // so that after this call, no text node spans a newline. Element
    // children are recursed into so a `<span>foo\nbar</span>` becomes
    // `<span>foo</span>\n<span>bar</span>` - necessary so we can wrap
    // whole lines cleanly.
    function splitTextNodesOnNewlines(parent) {
      var child = parent.firstChild;
      while (child) {
        var next = child.nextSibling;
        if (child.nodeType === Node.TEXT_NODE) {
          var v = child.nodeValue;
          var nlIdx = v.indexOf('\n');
          if (nlIdx !== -1) {
            // Split into: [pre-\n text node] [\n text node] [rest text node]
            // (the rest is handled on the next iteration).
            var rest = v.substring(nlIdx + 1);
            if (nlIdx > 0) {
              parent.insertBefore(document.createTextNode(v.substring(0, nlIdx)), child);
            }
            parent.insertBefore(document.createTextNode('\n'), child);
            if (rest.length > 0) {
              child.nodeValue = rest;
              // Re-process the shortened node in case it still contains newlines.
              next = child;
            } else {
              parent.removeChild(child);
            }
          }
        } else if (child.nodeType === Node.ELEMENT_NODE) {
          // If this element spans a newline, splitting DOM would fragment
          // the syntax spans. In practice hljs's Go grammar doesn't emit
          // spans that cross diff-line boundaries because our lines never
          // include multiple statements, so this branch is defensive: if
          // it does happen, split the element into siblings.
          if ((child.textContent || '').indexOf('\n') !== -1) {
            splitElementOnNewlines(child);
            next = child.nextSibling;
          }
        }
        child = next;
      }
    }

    function splitElementOnNewlines(el) {
      // Walk children; on the first newline text node, cut the element
      // into two: everything before stays in `el`, everything after goes
      // into a clone that becomes el's next sibling (with the newline as
      // a bare text node between them).
      var parent = el.parentNode;
      var child = el.firstChild;
      while (child) {
        var next = child.nextSibling;
        if (child.nodeType === Node.TEXT_NODE && child.nodeValue.indexOf('\n') !== -1) {
          var v = child.nodeValue;
          var nlIdx = v.indexOf('\n');
          var pre = v.substring(0, nlIdx);
          var post = v.substring(nlIdx + 1);
          if (pre.length > 0) {
            child.nodeValue = pre;
          } else {
            el.removeChild(child);
          }
          // Insert bare newline text node after el.
          var newline = document.createTextNode('\n');
          parent.insertBefore(newline, el.nextSibling);
          // Move remaining siblings (and the post-newline chunk) into a
          // clone of el placed after the newline.
          var clone = el.cloneNode(false);
          parent.insertBefore(clone, newline.nextSibling);
          if (post.length > 0) {
            clone.appendChild(document.createTextNode(post));
          }
          var sib = next;
          while (sib) {
            var s2 = sib.nextSibling;
            clone.appendChild(sib);
            sib = s2;
          }
          // Recurse: clone itself may still contain newlines.
          if ((clone.textContent || '').indexOf('\n') !== -1) {
            splitElementOnNewlines(clone);
          }
          return;
        }
        child = next;
      }
    }
  }

  register();
})();
