// d3 jelek mending react flow

// "use client";
// import { useState } from "react";
// import D3Tree, { TreeNode } from "@/components/d3tree";

// const Page = () => {
//   const [treeData, setTreeData] = useState<TreeNode>({
//     name: "root",
//     children: [],
//   });

//   const [parentName, setParentName] = useState("root");
//   const [childName, setChildName] = useState("");

//   const findAndAddChild = (
//     node: TreeNode,
//     targetName: string,
//     newChild: TreeNode
//   ): boolean => {
//     if (node.name === targetName) {
//       if (!node.children) node.children = [];
//       node.children.push(newChild);
//       return true;
//     }
//     if (node.children) {
//       for (const child of node.children) {
//         if (findAndAddChild(child, targetName, newChild)) return true;
//       }
//     }
//     return false;
//   };

//   const handleAddChild = () => {
//     if (!childName.trim()) return;

//     // Clone tree
//     const newTree = structuredClone(treeData);
//     const success = findAndAddChild(newTree, parentName, {
//       name: childName.trim(),
//     });

//     if (success) {
//       setTreeData(newTree);
//       setChildName("");
//     } else {
//       alert("Parent not found.");
//     }
//   };

//   const collectNodeNames = (node: TreeNode): string[] => {
//     const names = [node.name];
//     if (node.children) {
//       for (const child of node.children) {
//         names.push(...collectNodeNames(child));
//       }
//     }
//     return names;
//   };

//   const allNodeNames = collectNodeNames(treeData);

//   return (
//     <div className="w-full max-h-screen flex flex-col items-center justify-start p-4 gap-4">
//       <div className="bg-gray-100 p-4 rounded shadow w-full max-w-md">
//         <h2 className="font-bold mb-2">Add Child Node</h2>
//         <div className="flex flex-col gap-2">
//           <label>
//             Parent:
//             <select
//               className="w-full border p-1"
//               value={parentName}
//               onChange={(e) => setParentName(e.target.value)}
//             >
//               {allNodeNames.map((name) => (
//                 <option key={name} value={name}>
//                   {name}
//                 </option>
//               ))}
//             </select>
//           </label>
//           <label>
//             New Child Name:
//             <input
//               className="w-full border p-1"
//               type="text"
//               value={childName}
//               onChange={(e) => setChildName(e.target.value)}
//               placeholder="e.g., childX"
//             />
//           </label>
//           <button
//             className="bg-blue-500 text-white px-3 py-1 rounded"
//             onClick={handleAddChild}
//           >
//             Add Child
//           </button>
//         </div>
//       </div>

//       <div className="w-50%  border mt-4">
//         <D3Tree data={treeData} />
//       </div>
//     </div>
//   );
// };

// export default Page;
// "use client";

// import dynamic from "next/dynamic";
// import { useState, useRef, useEffect } from "react";

// const Tree = dynamic(() => import("react-d3-tree"), { ssr: false });

// const treeData =
//   // [
//   //   // Recipe 1: Brick = Mud + Fire
//   //   {
//   //     name: "Brick",
//   //     children: [
//   //       {
//   //         name: "Mud",
//   //         children: [{ name: "Water" }, { name: "Earth" }],
//   //       },
//   //       { name: "Fire" },
//   //     ],
//   //   },
//   //   // Recipe 2: Brick = Clay + Stone
//   //   {
//   //     name: "Brick",
//   //     children: [
//   //       {
//   //         name: "Clay",
//   //         children: [
//   //           {
//   //             name: "Mud",
//   //             children: [{ name: "Water" }, { name: "Earth" }],
//   //           },
//   //           {
//   //             name: "Sand",
//   //             children: [
//   //               {
//   //                 name: "Stone",
//   //                 children: [
//   //                   {
//   //                     name: "Lava",
//   //                     children: [{ name: "Earth" }, { name: "Fire" }],
//   //                   },
//   //                   { name: "Air" },
//   //                 ],
//   //               },
//   //             ],
//   //           },
//   //         ],
//   //       },
//   //       {
//   //         name: "Stone",
//   //         children: [
//   //           {
//   //             name: "Lava",
//   //             children: [{ name: "Earth" }, { name: "Fire" }],
//   //           },
//   //           { name: "Air" },
//   //         ],
//   //       },
//   //     ],
//   //   },
//   // ];

//   [
//     {
//       name: "Brick (All Recipes)",
//       children: [
//         {
//           name: "Brick",
//           children: [
//             {
//               name: "Mud",
//               children: [{ name: "Water" }, { name: "Earth" }],
//             },
//             { name: "Fire" },
//           ],
//         },
//         {
//           name: "Brick",
//           children: [
//             {
//               name: "Clay",
//               children: [
//                 {
//                   name: "Mud",
//                   children: [{ name: "Water" }, { name: "Earth" }],
//                 },
//                 {
//                   name: "Sand",
//                   children: [
//                     {
//                       name: "Stone",
//                       children: [
//                         {
//                           name: "Lava",
//                           children: [{ name: "Earth" }, { name: "Fire" }],
//                         },
//                         { name: "Air" },
//                       ],
//                     },
//                   ],
//                 },
//               ],
//             },
//             {
//               name: "Stone",
//               children: [
//                 {
//                   name: "Lava",
//                   children: [{ name: "Earth" }, { name: "Fire" }],
//                 },
//                 { name: "Air" },
//               ],
//             },
//           ],
//         },
//       ],
//     },
//   ];
// const containerStyles = {
//   width: "100%",
//   height: "100vh",
// };
// function TreeDiagram() {
//   const [translate, setTranslate] = useState({ x: 400, y: 100 });
//   const containerRef = useRef(null);

//   useEffect(() => {
//     if (containerRef.current) {
//       const { width, height } = containerRef.current.getBoundingClientRect();
//       setTranslate({ x: width / 2, y: 100 });
//     }
//   }, []);

//   return (
//     <div style={{ width: "100%", height: "100vh" }} ref={containerRef}>
//       <Tree
//         data={treeData}
//         translate={translate}
//         orientation="vertical"
//         pathFunc="elbow"
//         zoomable
//         collapsible={false}
//       />
//     </div>
//   );
// }

// export default function Page() {
//   return (
//     <main>
//       <TreeDiagram />
//     </main>
//   );
// }
