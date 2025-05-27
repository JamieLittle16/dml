// Package latex provides LaTeX template functionality for DML
package latex

// TexTemplate is the LaTeX document template for rendering math expressions
const TexTemplate = `\documentclass[border=2pt,preview]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{amsfonts}
\usepackage{mathtools}
\usepackage[dvipsnames,svgnames,table]{xcolor}
%s
\begin{document}
\pagecolor{%s}
\color{%s}
%s
\end{document}`

// FullDocTemplate is the LaTeX document template for rendering entire documents
const FullDocTemplate = `\documentclass[border=3pt,preview]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{amsfonts}
\usepackage{mathtools}
\usepackage[dvipsnames,svgnames,table]{xcolor}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
\usepackage{verbatim}
%s
\begin{document}
\pagecolor{%s}
\color{%s}
%s
\end{document}`