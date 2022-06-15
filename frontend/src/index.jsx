import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Index } from "./routes/Index";
import { Category } from "./routes/Category";
import { Latest } from "./routes/Latest";
import { createGlobalStyle } from "styled-components";
import { Toaster } from "react-hot-toast";

const GlobalStyle = createGlobalStyle`
  body {
    margin: 0;
    padding: 0;
    font-family: "Inter", sans-serif;
    background-color: #fff;
    font-size: 1.1rem;
    line-height: 1.65rem;
  }

  a {
    color: black;
  }

  a:visited {
    color: #7e4fe3;
  }

  h2 {
    margin: 0.1rem 0 1rem 0;
  }

  main {
    min-height: 75vh;
  }
`;

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <React.StrictMode>
    <GlobalStyle />
    <Toaster position="top-right" reverseOrder={false} />
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />}>
          <Route index element={<Latest />} />
          <Route path="category/:categoryId" element={<Category />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
