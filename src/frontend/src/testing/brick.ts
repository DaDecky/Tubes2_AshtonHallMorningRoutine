type recipeNode = {
  name: string;
  recipes?: [recipeNode, recipeNode][];
};
export const data: recipeNode = {
  name: "Brick",
  recipes: [
    [
      {
        name: "Mud",
        recipes: [
          [
            {
              name: "Water",
            },
            {
              name: "Earth",
            },
          ],
        ],
      },
      {
        name: "Fire",
      },
    ],
    [
      {
        name: "Clay",
        recipes: [
          [
            {
              name: "Mud",
              recipes: [
                [
                  {
                    name: "Water",
                  },
                  {
                    name: "Earth",
                  },
                ],
              ],
            },
            {
              name: "Sand",
              recipes: [
                [
                  {
                    name: "Stone",
                    recipes: [
                      [
                        {
                          name: "Lava",
                          recipes: [
                            [
                              {
                                name: "Earth",
                              },
                              {
                                name: "Fire",
                              },
                            ],
                          ],
                        },
                        {
                          name: "Air",
                        },
                      ],
                      [
                        {
                          name: "Earth",
                        },
                        {
                          name: "Pressure",
                          recipes: [[{ name: "Air" }, { name: "Air" }]],
                        },
                      ],
                    ],
                  },
                  {
                    name: "Air",
                  },
                ],
              ],
            },
          ],
        ],
      },
      {
        name: "Stone",
        recipes: [
          [
            {
              name: "Lava",
              recipes: [
                [
                  {
                    name: "Earth",
                  },
                  {
                    name: "Fire",
                  },
                ],
              ],
            },
            {
              name: "Air",
            },
          ],
          [
            {
              name: "Earth",
            },
            {
              name: "Pressure",
              recipes: [[{ name: "Air" }, { name: "Air" }]],
            },
          ],
        ],
      },
    ],
  ],
};
