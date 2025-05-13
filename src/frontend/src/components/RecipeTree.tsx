"use client";

import ReactFlow, {
  Background,
  Controls,
  MiniMap,
  Position,
  Node,
  Edge,
  MarkerType,
} from "reactflow";
import "reactflow/dist/style.css";
import { getLayoutedElements } from "@/utils/graphLayout";
import Image from "next/image";
import { useEffect, useState, useRef } from "react";

type RecipeNode = {
  name: string;
  recipes?: [RecipeNode, RecipeNode][];
};

type recipeTreeProps = {
  data: RecipeNode | undefined;
  direction?: "TB" | "LR";
  useImage?: boolean;
  algorithm?: "BFS" | "DFS";
  liveplay?: boolean;
  speed?: number;
};

const nodeDefaults = {
  sourcePosition: Position.Bottom,
  targetPosition: Position.Top,
};

const comboStyle = {
  width: 10,
  height: 10,
  background: "transparent",
  border: "none",
};

const generateGraphElementsAnimated = (
  data: RecipeNode,
  direction: "TB" | "LR" = "TB",
  algorithm: "BFS" | "DFS" = "BFS"
) => {
  const allNodes: Node[] = [];
  const allEdges: Edge[] = [];
  let nodeIdCounter = 1;
  let comboCounter = 1;
  let recipeCounter = 1;

  if (algorithm === "BFS") {
    const queue: { node: RecipeNode; parentId?: string; isRoot?: boolean }[] = [];
    queue.push({ node: data, isRoot: true });

    while (queue.length > 0) {
      const current = queue.shift()!;
      const nodeId = `node${nodeIdCounter++}`;

      allNodes.push({
        id: nodeId,
        data: {
          label: (
            <div className="flex items-center">
              <div className="w-12 h-12 relative">
                <Image
                  src={`/images/${current.node.name}.svg`}
                  alt={current.node.name}
                  fill
                />
              </div>
              <span className="text-xl">{current.node.name}</span>
            </div>
          ),
        },
        ...nodeDefaults,
      } as Node);

      if (current.node.recipes && current.node.recipes.length > 0) {
        for (const [left, right] of current.node.recipes) {
          const comboId = `combo${comboCounter++}`;

          // Add combo node
          allNodes.push({
            id: comboId,
            data: { label: "" },
            style: comboStyle,
            ...nodeDefaults,
          } as Node);

          // Add edges
          allEdges.push({
            id: `e-${comboId}-${nodeId}`,
            source: comboId,
            target: nodeId,
            label: current.isRoot ? `Recipe ${recipeCounter++}` : undefined,
            style: { stroke: "#4caf50", strokeWidth: 2 },
            markerEnd: { type: MarkerType.ArrowClosed },
          });

          // Process children (BFS)
          queue.push({ node: left, parentId: comboId });
          queue.push({ node: right, parentId: comboId });
        }
      }

      // If this node has a parent, connect it
      if (current.parentId) {
        allEdges.push({
          id: `e-${nodeId}-${current.parentId}`,
          source: nodeId,
          target: current.parentId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });
      }
    }
  } else {
    const stack: { node: RecipeNode; parentId?: string; isRoot?: boolean }[] = [];
    stack.push({ node: data, isRoot: true });

    while (stack.length > 0) {
      const current = stack.pop()!;
      const nodeId = `node${nodeIdCounter++}`;

      allNodes.push({
        id: nodeId,
        data: {
          label: (
            <div className="flex items-center">
              <div className="w-12 h-12 relative">
                <Image
                  src={`/images/${current.node.name}.svg`}
                  alt={current.node.name}
                  fill
                />
              </div>
              <span className="text-xl">{current.node.name}</span>
            </div>
          ),
        },
        ...nodeDefaults,
      } as Node);

      if (current.node.recipes && current.node.recipes.length > 0) {
        const recipes = [...current.node.recipes].reverse();
        
        for (const [left, right] of recipes) {
          const comboId = `combo${comboCounter++}`;

          allNodes.push({
            id: comboId,
            data: { label: "" },
            style: comboStyle,
            ...nodeDefaults,
          } as Node);

          allEdges.push({
            id: `e-${comboId}-${nodeId}`,
            source: comboId,
            target: nodeId,
            label: current.isRoot ? `Recipe ${recipeCounter++}` : undefined,
            style: { stroke: "#4caf50", strokeWidth: 2 },
            markerEnd: { type: MarkerType.ArrowClosed },
          });

          stack.push({ node: right, parentId: comboId });
          stack.push({ node: left, parentId: comboId });
        }
      }

      if (current.parentId) {
        allEdges.push({
          id: `e-${nodeId}-${current.parentId}`,
          source: nodeId,
          target: current.parentId,
          label: "combine",
          style: { strokeDasharray: "5 5" },
        });
      }
    }
  }

  return getLayoutedElements(allNodes, allEdges, direction);
};

export default function RecipeTree({
  data,
  direction = "TB",
  algorithm = "BFS",
  liveplay = false,
  speed = 500,
}: recipeTreeProps) {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const speedRef = useRef(speed);
  const liveplayRef = useRef(liveplay);


  useEffect(() => {
    speedRef.current = speed;
    liveplayRef.current = liveplay
  }, [speed, liveplay]);

  useEffect(() => {
    if (!data) return;
    const { nodes: allNodes, edges: allEdges } = generateGraphElementsAnimated(
      data,
      direction,
      algorithm
    );

    setNodes([]);
    setEdges([]);
    
    if (liveplayRef.current) {
      const nodeInterval = setInterval(() => {
        setNodes((prevNodes) => {
          if (allNodes.length > prevNodes.length) {
            const newNodes = [...prevNodes, allNodes[prevNodes.length]];
  
            setEdges((prevEdges) => {
              if (allEdges.length > prevEdges.length) {
                return [...prevEdges, allEdges[prevEdges.length]];
              }
              return prevEdges;
            });
  
            return newNodes;
          } else {
            clearInterval(nodeInterval);
            return prevNodes;
          }
        });
      }, speedRef.current);
      return () => clearInterval(nodeInterval);
    }
    else {
      setNodes(allNodes);
      setEdges(allEdges);
    }

    
  }, [data, direction, algorithm]);

  
  return (
    <div style={{ width: "100%", height: "100%" }}>
      
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        defaultEdgeOptions={{ style: { strokeWidth: 2 } }}
      >
        <Background />
        <Controls />
        <MiniMap />
      </ReactFlow>
    </div>
  );
}