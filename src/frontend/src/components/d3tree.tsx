// d3 jelek pake react flow aja

// "use client";
// import * as d3h from "d3-hierarchy";
// import * as d3 from "d3";
// import { useEffect, useRef } from "react";
// export type TreeNode = {
//   name: string;
//   children?: TreeNode[];
// };

// export default function D3Tree({
//   data,
//   nodeWidth = 120,
//   nodeHeight = 60,
//   marginTop = 20,
//   marginRight = 20,
//   marginBottom = 20,
//   marginLeft = 20,
// }: {
//   data: TreeNode;
//   nodeWidth?: number;
//   nodeHeight?: number;
//   marginTop?: number;
//   marginRight?: number;
//   marginBottom?: number;
//   marginLeft?: number;
// }) {
//   const svgRef = useRef<SVGSVGElement>(null);
//   const gRef = useRef<SVGGElement>(null);

//   const root = d3h.hierarchy(data);
//   const treeLayout = d3h.tree<TreeNode>().nodeSize([nodeHeight, nodeWidth]);
//   treeLayout(root);

//   const nodes = root.descendants();
//   const links = root.links();

//   const xExtent = d3.extent(nodes, (d) => d.x) as [number, number];
//   const yExtent = d3.extent(nodes, (d) => d.y) as [number, number];

//   const width = yExtent[1] - yExtent[0] + marginLeft + marginRight;
//   const height = xExtent[1] - xExtent[0] + marginTop + marginBottom;

//   const translateX = marginLeft - yExtent[0];
//   const translateY = marginTop - xExtent[0];

//   // Enable zoom and pan
//   useEffect(() => {
//     if (!svgRef.current || !gRef.current) return;

//     const svg = d3.select(svgRef.current);
//     const g = d3.select(gRef.current);

//     const zoom = d3
//       .zoom<SVGSVGElement, unknown>()
//       .scaleExtent([0.2, 3]) // Min and max zoom
//       .on("zoom", (event) => {
//         g.attr("transform", event.transform);
//       });

//     svg.call(zoom);
//   }, []);

//   return (
//     <div className="overflow-hidden">
//       <svg
//         ref={svgRef}
//         width="100%"
//         height="100%"
//         viewBox={`0 0 ${width} ${height}`}
//         className="border border-gray-300 bg-white"
//       >
//         <g ref={gRef} transform={`translate(${translateX}, ${translateY})`}>
//           {links.map((link, i) => (
//             <line
//               key={i}
//               x1={link.source.y}
//               y1={link.source.x}
//               x2={link.target.y}
//               y2={link.target.x}
//               stroke="black"
//             />
//           ))}
//           {nodes.map((node, i) => (
//             <g key={i} transform={`translate(${node.y}, ${node.x})`}>
//               <circle r="4" fill="black" />
//               <text dx="6" dy="4" fontSize="12px">
//                 {node.data.name}
//               </text>
//             </g>
//           ))}
//         </g>
//       </svg>
//     </div>
//   );
// }
