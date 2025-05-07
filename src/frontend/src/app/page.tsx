import Link from "next/link";
import React from "react";

const page = () => {
  return (
    <div>
      <Link
        className="bg-blue-500 text-white font-bold py-2 px-4 rounded hover:bg-blue-700"
        href={"/dummy"}
      >
        Dummy Live Update Graph
      </Link>
    </div>
  );
};

export default page;
