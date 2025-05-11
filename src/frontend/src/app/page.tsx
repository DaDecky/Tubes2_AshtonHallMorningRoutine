// "use client";

// import Link from "next/link";
// import React, { useState } from "react";
// import { Label } from "@/components/ui/label";
// import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
// import { redirect } from "next/navigation";

// export enum Mode {
//   SHORTEST = "Shortest",
//   MULTI = "Multi",
// }
// export enum Algorithm {
//   BFS = "BFS",
//   DFS = "DFS",
// }

// const Page = () => {
//   const [mode, setMode] = useState<Mode>(Mode.SHORTEST);
//   const [algorithm, setAlgorithm] = useState<Algorithm>(Algorithm.BFS);
//   const [maxRecipe, setMaxRecipe] = useState<number>(-1);

//   const handleQuery = async () => {
//     const response = await fetch(
//       `/api/query?algo=${algorithm}&mode=${mode}&max=${maxRecipe}`
//     );

//     const data = await response.json();
//     console.log("Query response:", data);
//   };

//   return (
//     <div className="flex flex-col items-center  min-h-screen bg-gray-100">
//       <Link
//         className="bg-blue-500 text-white font-bold py-2 px-4 rounded hover:bg-blue-700"
//         href={"/dummy"}
//       >
//         Dummy Live Update Graph
//       </Link>

//       <div>
//         <div>Choose Algorithm: </div>
//         <RadioGroup className="flex" defaultValue={algorithm as string}>
//           {Object.values(Algorithm).map((algo) => (
//             <div className="flex items-center space-x-2" key={algo}>
//               <RadioGroupItem
//                 value={algo}
//                 id={algo}
//                 onClick={() => setAlgorithm(algo)}
//               />
//               <Label htmlFor={algo}>{algo}</Label>
//             </div>
//           ))}
//         </RadioGroup>

//         <div>Choose Mode : </div>
//         <RadioGroup className="flex" defaultValue={mode as string}>
//           {Object.values(Mode).map((mode) => (
//             <div className="flex items-center space-x-2" key={mode}>
//               <RadioGroupItem
//                 value={mode}
//                 id={mode}
//                 onClick={() => setMode(mode)}
//               />
//               <Label htmlFor={mode}>{mode}</Label>
//             </div>
//           ))}
//         </RadioGroup>

//         {mode === Mode.MULTI && (
//           <div>
//             <label htmlFor="max-recipe">Max Recipe: </label>
//             <input
//               type="number"
//               id="max-recipe"
//               value={maxRecipe}
//               onChange={(e) => setMaxRecipe(Number(e.target.value))}
//               className="border border-gray-300 rounded p-2"
//             />
//           </div>
//         )}
//       </div>

//       <button
//         className="bg-blue-500 text-white font-bold py-2 px-4 rounded hover:bg-blue-700 mt-4 hover:cursor-pointer"
//         onClick={handleQuery}
//       >
//         Query Result
//       </button>
//       <p className="mt-4">Current mode: {mode}</p>
//       <p className="mt-4">Current algorithm: {algorithm}</p>
//       <p className="mt-4">
//         {mode === Mode.MULTI ? `Max Recipe: ${maxRecipe}` : ""}
//       </p>
//     </div>
//   );
// };

// export default Page;

import { redirect } from "next/navigation";
import React from "react";

const page = () => {
  redirect("/rflow");
  return <div>page</div>;
};

export default page;
