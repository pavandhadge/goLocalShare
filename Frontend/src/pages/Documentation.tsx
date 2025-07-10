import * as React from "react";
import ReactMarkdown from "react-markdown";
import { Card } from "@/components/ui/card";
import rehypeRaw from "rehype-raw";
import rehypeSanitize from "rehype-sanitize";
import remarkGfm from "remark-gfm";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark as githubTheme } from "react-syntax-highlighter/dist/esm/styles/prism";
import type { ReactMarkdownProps } from "react-markdown/lib/ast-to-react";
import { Link } from "react-router-dom";
const DOCS = [
  { label: "Introduction", file: "introduction.md" },
  { label: "Quick Start", file: "quickstart.md" },
  { label: "Installation", file: "installation.md" },
  { label: "Usage & CLI", file: "usage.md" },
  { label: "Architecture", file: "architecture.md" },
  { label: "Security", file: "security.md" },
  { label: "Cloud Upload", file: "cloud.md" },
  { label: "API & Web UI", file: "api.md" },
  { label: "Configuration", file: "configuration.md" },
  { label: "Troubleshooting", file: "troubleshooting.md" },
  { label: "Contributing", file: "contributing.md" },
  { label: "Changelog", file: "changelog.md" },
];

const DOCS_PATH = "/docs/";

export const DocumentationButton = ({ className = "", ...props }) => (
  <Link
    to="/documentation"
    className={`inline-flex border border-white items-center gap-2 px-4 py-2 rounded-lg font-semibold bg-white text-black hover:bg-blue-600 shadow transition-colors duration-150 ${className}`}
    {...props}
  >
    <svg width="20" height="20" fill="none" viewBox="0 0 24 24"><path stroke="currentColor" strokeWidth="2" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Zm0 0V2m10 2V2"/></svg>
    Documentation
  </Link>
);

const Documentation: React.FC = () => {
  const [selected, setSelected] = React.useState(0);
  const [content, setContent] = React.useState<string>("");
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  React.useEffect(() => {
    fetch(DOCS_PATH + DOCS[selected].file)
      .then((res) => res.text())
      .then(setContent);
  }, [selected]);

  // Close sidebar on doc select (mobile)
  React.useEffect(() => { setSidebarOpen(false); }, [selected]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-blue-100 flex flex-col">
      <header className="bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-lg">
        <div className="container mx-auto px-4 py-10 text-center">
          <h1 className="text-4xl md:text-5xl font-extrabold mb-2 tracking-tight drop-shadow-sm">
            Documentation
          </h1>
          <p className="text-lg text-blue-100 font-medium">
            Deep, user-friendly docs for goLocalShare
          </p>
        </div>
      </header>
      <main className="flex flex-1 flex-col md:flex-row container mx-auto px-2 sm:px-4 py-4 sm:py-8 gap-4 md:gap-8 w-full">
        {/* Mobile sidebar toggle */}
        <button
          className="md:hidden mb-4 self-start px-4 py-2 rounded-lg bg-blue-600 text-white font-semibold shadow hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
          onClick={() => setSidebarOpen((v) => !v)}
        >
          {sidebarOpen ? "Hide Sections" : "Show Sections"}
        </button>
        {/* Sidebar */}
        <nav
          className={`z-20 md:z-auto bg-white md:bg-transparent rounded-lg md:rounded-none shadow md:shadow-none border md:border-0 w-full md:w-64 min-w-0 md:min-w-[180px] md:max-w-[220px] transition-all duration-200 ${
            sidebarOpen ? "block" : "hidden"
          } md:block`}
        >
          <ul className="md:sticky md:top-24 space-y-2 p-2 md:p-0">
            {DOCS.map((doc, idx) => (
              <li key={doc.file}>
                <button
                  className={`w-full text-left px-4 py-2 rounded-lg font-medium transition-colors duration-150 ${
                    selected === idx
                      ? "bg-blue-600 text-white shadow"
                      : "bg-white text-blue-700 hover:bg-blue-100"
                  }`}
                  onClick={() => setSelected(idx)}
                >
                  {doc.label}
                </button>
              </li>
            ))}
          </ul>
        </nav>
        {/* Content */}
        <section className="flex-1 min-w-0">
          <Card className="p-4 sm:p-6 md:p-10 bg-[#fcfcfd] shadow-xl border border-gray-200 max-w-none prose prose-base sm:prose-lg prose-headings:text-gray-900 prose-headings:font-bold prose-headings:tracking-tight prose-headings:mt-8 sm:prose-headings:mt-10 prose-headings:mb-4 sm:prose-headings:mb-5 prose-h1:text-3xl sm:prose-h1:text-4xl prose-h2:text-2xl sm:prose-h2:text-3xl prose-h3:text-xl sm:prose-h3:text-2xl prose-h4:text-lg sm:prose-h4:text-xl prose-h5:text-base sm:prose-h5:text-lg prose-h6:text-sm sm:prose-h6:text-base prose-blockquote:bg-gray-50 prose-blockquote:border-l-4 prose-blockquote:border-gray-300 prose-blockquote:pl-4 sm:prose-blockquote:pl-5 prose-blockquote:italic prose-blockquote:text-gray-700 prose-code:bg-gray-100 prose-code:text-gray-900 prose-code:rounded prose-code:px-1 prose-code:py-0.5 prose-pre:bg-gray-100 prose-pre:text-gray-900 prose-pre:rounded-lg prose-pre:p-3 sm:prose-pre:p-5 prose-pre:overflow-x-auto prose-table:rounded-lg prose-table:overflow-hidden prose-table:border prose-table:border-gray-200 prose-th:bg-gray-50 prose-th:text-gray-900 prose-td:bg-white prose-td:border-gray-100 prose-td:p-2 sm:prose-td:p-4 prose-th:p-2 sm:prose-th:p-4 prose-table:my-4 sm:prose-table:my-6 prose-table:w-full prose-tr:hover:bg-gray-50 prose-p:leading-relaxed prose-p:my-3 sm:prose-p:my-5 prose-li:my-1 sm:prose-li:my-2 prose-ul:pl-5 sm:prose-ul:pl-7 prose-ol:pl-5 sm:prose-ol:pl-7 prose-hr:my-6 sm:prose-hr:my-10 prose-a:text-blue-700 prose-a:underline prose-a:font-medium prose-a:hover:text-blue-900 text-gray-900 font-[Segoe_UI,Liberation_Sans,Arial,sans-serif]">
            <ReactMarkdown
              remarkPlugins={[remarkGfm]}
              rehypePlugins={[rehypeRaw, rehypeSanitize]}
              components={{
                a: ({ node, ...props }) => (
                  <a
                    {...props}
                    className="text-blue-600 underline hover:text-blue-800"
                    target="_blank"
                    rel="noopener noreferrer"
                  />
                ),
                table: ({ node, ...props }) => (
                  <table
                    {...props}
                    className="w-full border border-gray-200 rounded-lg overflow-hidden my-4"
                  />
                ),
                th: ({ node, ...props }) => (
                  <th
                    {...props}
                    className="bg-gray-100 text-gray-700 p-2 sm:p-3 border border-gray-200"
                  />
                ),
                td: ({ node, ...props }) => (
                  <td
                    {...props}
                    className="bg-white p-2 sm:p-3 border border-gray-100"
                  />
                ),
                blockquote: ({ node, ...props }) => (
                  <blockquote
                    {...props}
                    className="bg-gray-50 border-l-4 border-gray-300 pl-3 sm:pl-4 italic text-gray-600 my-3 sm:my-4"
                  />
                ),
                code({ node, inline, className, children, ...props }: { node: any; inline?: boolean; className?: string; children: React.ReactNode }) {
                  const match = /language-(\w+)/.exec(className || "");
                  return !inline ? (
                    <SyntaxHighlighter
                      style={githubTheme}
                      language={match ? match[1] : undefined}
                      PreTag="div"
                      customStyle={{
                        background: "#f6f8fa",
                        borderRadius: "0.5rem",
                        fontSize: "0.95em",
                        margin: "1.2em 0",
                        padding: "1.2em",
                      }}
                      {...props}
                    >
                      {String(children).replace(/\n$/, "")}
                    </SyntaxHighlighter>
                  ) : (
                    <code
                      className="bg-gray-100 rounded px-1 py-0.5 text-sm font-mono"
                      {...props}
                    >
                      {children}
                    </code>
                  );
                },
              }}
            >
              {content}
            </ReactMarkdown>
          </Card>
        </section>
      </main>
    </div>
  );
};

export default Documentation; 