
import { Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import Documentation from "./pages/Documentation";


function App() {
  return (
    <Routes>
      <Route path="/" element={<Index />} />
      <Route path="/documentation" element={<Documentation />} />
    </Routes>
  );
}

export default App;
