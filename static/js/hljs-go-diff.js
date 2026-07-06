// Registers a `go-diff` language on highlight.js that renders a unified diff
// with Go syntax highlighting inside every +/-/context line. Depends on
// highlight.js already having the built-in `go` and `diff` languages loaded.
(function () {
  function register() {
    if (!window.hljs || typeof window.hljs.registerLanguage !== 'function') {
      // code.min.js hasn't finished evaluating yet; try again on the next tick.
      return setTimeout(register, 0);
    }

    window.hljs.registerLanguage('go-diff', function (hljs) {
      return {
        name: 'Go Diff',
        aliases: ['diff-go', 'godiff'],
        contains: [
          {
            // Hunk headers: @@ -1,2 +3,4 @@
            className: 'meta',
            relevance: 10,
            match: /^(@@\s+-\d+(?:,\d+)?\s+\+\d+(?:,\d+)?\s+@@.*)$/m,
          },
          {
            // File headers and index lines.
            className: 'comment',
            variants: [
              { begin: /^(---|\+\+\+|\*{3}|Index:|diff --git|={3,}).*$/m },
              { match: /^\*{15}$/m },
            ],
          },
          {
            className: 'addition',
            begin: /^\+/,
            end: /$/,
            subLanguage: 'go',
          },
          {
            className: 'deletion',
            begin: /^-/,
            end: /$/,
            subLanguage: 'go',
          },
          {
            // Everything else on a line — context lines with a leading space,
            // and any un-prefixed lines the author left as-is. No diff
            // background, but the Go sub-language still colours the tokens.
            begin: /^[^+\-@]/,
            end: /$/,
            subLanguage: 'go',
          },
        ],
      };
    });
  }

  register();
})();
