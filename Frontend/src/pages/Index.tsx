
import { Button } from "@/components/ui/button";
import {
  ArrowRight,
  Shield,
  Cpu,
  Server,
  Share2,
  Clock,
  Lock,
  FileText,
  Database,
  Github,
  Download,
} from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DocumentationButton } from "./Documentation";

const Index: React.FC = () => {
  // GitHub username to be prominently displayed
  const githubUsername = "pavandhadge/goLocalShare";
  // Session duration (static for now)
  const sessionDuration = "1 hour";
  // Token system description
  const tokenDescription = "Access is protected by a time-limited token. Only users with the correct token can browse or download files during the session.";

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-blue-100">
      {/* Hero Section with professional color scheme */}
      <header className="bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-lg max-h-[75vh]">
        <div className="container mx-auto px-4 py-24 md:py-32">
          <div className="max-w-3xl mx-auto text-center">
            <div className="flex justify-center mb-8">
              {/* Go Gopher Mascot */}
              <img
                src="/go-gopher-svgrepo-com.svg"
                alt="Go Gopher Mascot"
                className="h-32 w-32 animate-bounce drop-shadow-xl"
              />
            </div>
            <h1 className="text-5xl md:text-7xl font-extrabold mb-6 tracking-tight drop-shadow-sm">
              goLocalShare
            </h1>
            <p className="text-2xl md:text-3xl mb-10 text-blue-100 font-medium">
              Share files securely across your local network in seconds
            </p>
            <div className="flex flex-wrap justify-center gap-4">
              <Button
                onClick={() =>
                  window.open(
                    `https://github.com/pavandhadge/goFileShare/releases/tag/goFileSharev2`,
                    "_blank",
                  )
                }
                size="lg"
                className="gap-2 bg-white text-blue-600 hover:bg-blue-50 shadow-md"
              >
                <Download size={20} />
                Download
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="gap-2 border-white text-black hover:bg-white/10 shadow-md"
                onClick={() =>
                  window.open(`https://github.com/${githubUsername}`, "_blank")
                }
              >
                <Github size={20} />
                GitHub
              </Button>
              <DocumentationButton className="!px-4 !py-2 !rounded-lg !font-semibold !shadow-md !bg-blue-600 !text-white hover:!bg-blue-700" />
            </div>
          </div>
        </div>
        {/* Wave separator */}
        <div className="relative h-16 md:h-24 bg-gradient-to-r from-blue-600 to-blue-700">
          <svg
            className="absolute bottom-0 left-0 w-full h-full"
            viewBox="0 0 1440 100"
            preserveAspectRatio="none"
          >
            <path
              fill="rgb(249 250 251)"
              d="M0,0 C240,80 480,100 720,80 C960,60 1200,40 1440,80 L1440,100 L0,100 Z"
            ></path>
          </svg>
        </div>
      </header>

      {/* Features Section */}
      <section className="py-24 bg-white border-b border-blue-100">
        <div className="container mx-auto px-4">
          <h2 className="text-4xl md:text-5xl font-extrabold text-center mb-4 text-blue-700 tracking-tight">
            What It Does
          </h2>
          <p className="text-lg text-gray-600 text-center max-w-2xl mx-auto mb-16">
            Our tool makes sharing files simple and secure.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-10 p-4 rounded-2xl ">
            {features.map((feature, index) => (
              <Card
                key={index}
                className="border border-blue-200 shadow-lg hover:shadow-2xl transition-all duration-300 group overflow-hidden rounded-xl bg-white/80 backdrop-blur-md hover:bg-blue-100/70"
              >
                <CardHeader className="pb-2 relative">
                  <div className="absolute -right-8 -top-8 bg-blue-100 w-24 h-24 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></div>
                  <div className="mb-4 text-blue-500 relative z-10">
                    {feature.icon}
                  </div>
                  <CardTitle className="text-2xl font-bold relative z-10 text-blue-800">
                    {feature.title}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-foreground/80 text-base relative z-10">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* GitHub Focus Section - NEW */}
      <section className="py-20 bg-blue-50 border-b border-blue-100">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <Github className="h-16 w-16 mx-auto mb-6 text-blue-600" />
            <h2 className="text-3xl md:text-4xl font-bold mb-6 text-blue-700">Open Source on GitHub</h2>
            <p className="text-lg mb-8 text-blue-900">
              This project is completely open source. Check out the code,
              contribute, or report issues on GitHub.
            </p>
            <div className="bg-white rounded-lg shadow-lg p-6 mb-8 flex items-center justify-center">
              <div className="flex items-center bg-gray-100 px-5 py-3 rounded-md font-mono text-lg">
                <span className="text-gray-500 mr-2">github.com/</span>
                <span className="font-bold text-blue-600">
                  {githubUsername}
                </span>
              </div>
            </div>
            <Button
              size="lg"
              className="gap-2 bg-blue-600 hover:bg-blue-700 shadow-md"
              onClick={() =>
                window.open(`https://github.com/${githubUsername}`, "_blank")
              }
            >
              <Github size={20} />
              View on GitHub
            </Button>
          </div>
        </div>
      </section>

      {/* Usage Section */}
      <section className="py-24 bg-white border-b border-blue-100">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-2 gap-16 items-center">
            <div className="h-full flex flex-col justify-center">
              <h2 className="text-3xl md:text-4xl font-bold mb-6 text-blue-700">
                Simple to Use
              </h2>
              <p className="text-lg text-gray-600 mb-8">
                Start sharing files in seconds with our lightweight and powerful
                CLI tool.
              </p>
              <div className="space-y-4">
                <div className="bg-gray-900 text-gray-100 p-4 rounded-md font-mono text-sm shadow">
                  <p className="text-gray-400"># Share a single file for 2 hours</p>
                  <p>./goLocalShare --duration 2h /path/to/your/file.ext</p>
                </div>
                <div className="bg-gray-900 text-gray-100 p-4 rounded-md font-mono text-sm shadow">
                  <p className="text-gray-400"># Share a directory for 30 minutes</p>
                  <p>./goLocalShare --dir --duration 30m /path/to/your/directory</p>
                </div>
                <div className="bg-gray-900 text-gray-100 p-4 rounded-md font-mono text-sm shadow">
                  <p className="text-gray-400"># Upload a file to Cloudinary for 1 hour</p>
                  <p className="text-gray-400"># First time only: add --cloud-name, --cloud-key, --cloud-secret</p>
                  <p>./goLocalShare --cloud --duration 1h /path/to/your/file.ext</p>
                  <p className="text-gray-400"># Example (first time):</p>
                  <p>./goLocalShare --cloud --cloud-name &lt;name&gt; --cloud-key &lt;key&gt; --cloud-secret &lt;secret&gt; --duration 1h /path/to/your/file.ext</p>
                </div>
              </div>
            </div>
            <div className="h-full flex flex-col justify-center p-0">
  <div className="relative bg-gradient-to-br from-blue-100 to-blue-50 p-8 rounded-2xl shadow-lg h-full flex flex-col justify-center">
    <div className="bg-white/70 backdrop-blur-md p-6 rounded-xl shadow-inner">
      <h3 className="text-2xl font-bold mb-6 text-blue-700">Lightning Fast Setup</h3>
      <ol className="space-y-5">
        <li className="flex gap-4 items-start">
          <div className="bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center font-semibold">
            1
          </div>
          <p className="font-medium">
            Download the binary or build from source
          </p>
        </li>
        <li className="flex gap-4 items-start">
          <div className="bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center font-semibold">
            2
          </div>
          <p className="font-medium">
            Run the server command
          </p>
        </li>
        <li className="flex gap-4 items-start">
          <div className="bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center font-semibold">
            3
          </div>
          <p className="font-medium">
            Share the link + token with anyone on your network
          </p>
        </li>
      </ol>
    </div>
  </div>
</div>

          </div>
        </div>
      </section>

      {/* Technical Section with Gopher */}
      <section className="py-20 bg-gray-50 relative overflow-hidden">
        <div className="container mx-auto px-4 text-center relative z-10">
          <h2 className="text-3xl md:text-4xl font-bold mb-4 text-blue-700">
            Built for Security
          </h2>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto mb-12">
            Every aspect of goLocalShare is designed with security in mind
          </p>
          <div className="grid md:grid-cols-3 gap-6 max-w-4xl mx-auto">
            {securityFeatures.map((feature, index) => (
              <SecurityFeature
                key={index}
                icon={feature.icon}
                title={feature.title}
                description={feature.description}
              />
            ))}
          </div>
        </div>
        {/* Go Gopher in the background */}
        <div className="absolute -right-20 -bottom-20 opacity-10">
          <img
            src="https://go.dev/images/gophers/ladder.svg"
            alt="Go Gopher"
            className="w-64 h-64"
          />
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 bg-gradient-to-r from-blue-600 to-blue-700 text-white">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl md:text-4xl font-bold mb-6">
            Ready to start sharing?
          </h2>
          <p className="text-xl max-w-2xl mx-auto mb-10">
            Get goLocalShare today and share files with confidence
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <Button
              onClick={() =>
                window.open(
                  `https://github.com/pavandhadge/goFileShare/releases/tag/goFileSharev2`,
                  "_blank",
                )
              }
              size="lg"
              className="gap-2 bg-white text-blue-600 hover:bg-blue-50 shadow-md"
            >
              <Download size={20} />
              Download
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="gap-2 border-white text-black hover:bg-white/10 shadow-md"
              onClick={() =>
                window.open(`https://github.com/${githubUsername}`, "_blank")
              }
            >
              <Github size={20} />
              View on GitHub
            </Button>
            <DocumentationButton className="ml-2" />
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-300 py-12">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="flex items-center mb-6 md:mb-0">
              <Server className="h-8 w-8 mr-3 text-blue-500" />
              <span className="text-xl font-bold text-white">goLocalShare</span>
            </div>
            <div className="text-sm">
              <p>A simple, secure file sharing solution.</p>
              {/* <DocumentationButton className="ml-2" /> */}
            </div>
            <div className="flex items-center mt-6 md:mt-0">
              <Button
                variant="ghost"
                size="icon"
                className="text-gray-400 hover:text-black"
                onClick={() =>
                  window.open(`https://github.com/${githubUsername}`, "_blank")
                }
              >
                <Github size={48} />
              </Button>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

// Feature data
const features = [
  {
    icon: <Share2 className="h-12 w-12" />,
    title: "Easy File & Directory Sharing",
    description:
      "Share single files or entire directories over your local network with a single command.",
  },
  {
    icon: <Lock className="h-12 w-12" />,
    title: "Secure Token Authentication",
    description: "Access is protected by a time-limited token. Only users with the correct token can browse or download files.",
  },
  {
    icon: <Clock className="h-12 w-12" />,
    title: "Time-limited Session",
    description: "Session expires automatically after the configured duration (default: 1 hour).",
  },
  {
    icon: <Server className="h-12 w-12" />,
    title: "Cloud Upload & Auto-Delete",
    description: "Optionally upload files to Cloudinary for sharing. Files are automatically deleted from the cloud after the session ends.",
  },
  {
    icon: <Server className="h-12 w-12" />,
    title: "Security Headers",
    description: "CSP, XSS protection, no-sniff, and other secure HTTP headers are enforced.",
  },
  {
    icon: <FileText className="h-12 w-12" />,
    title: "Path & Symlink Protection",
    description: "Prevents directory traversal and symlink attacks. Only files within the shared path are accessible.",
  },
];

// Security features data
const securityFeatures = [
  {
    icon: <Lock />,
    title: "Authentication",
    description: "Time-limited tokens",
  },
  {
    icon: <Database />, title: "No Size Limit", description: "Share freely" },
  {
    icon: <Cpu />, title: "Lightweight", description: "Small footprint" },
];

// Helper component
const SecurityFeature = ({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) => (
  <div className="bg-white rounded-lg p-6 shadow hover:shadow-md transition-shadow text-center">
    <div className="bg-blue-100 text-blue-600 rounded-full p-3 w-14 h-14 mx-auto mb-4 flex items-center justify-center">
      {icon}
    </div>
    <h3 className="font-semibold text-lg mb-1">{title}</h3>
    <p className="text-gray-600 text-sm">{description}</p>
  </div>
);

export default Index;