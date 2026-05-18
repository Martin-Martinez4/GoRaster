import math
import bpy

# "maze"
magicNum = [109, 97, 122, 101]


global grid


grid = []

border_width = 0.5
border_width_2 = border_width * 2
width = 5
height = 5
tall = 2

BOTTOM = 4
RIGHT = 8
TOP = 1
LEFT = 2


def matrix_coords_to_array_coords(row, column, columns):
    return row * columns + column

def array_coords_to_matrix_coords(index, columns):
    # row, column
    return (math.floor(index / columns), index % columns)


def readMazeFile(filePath):


    with open(filePath, "rb") as f:
        # Check for magic number
        for i in range(4):
            byte = f.read(4)
            print(int.from_bytes(byte,  byteorder='little'))
            if (int.from_bytes(byte,  byteorder='little') != magicNum[i]):
                print("Magic Number Mismatch")
                return
            
        # Read columns and rows
        byte = f.read(4)
        columns = int.from_bytes(byte,  byteorder='little')

        byte = f.read(4)
        rows = int.from_bytes(byte,  byteorder='little')

            

        while True:
            byte = f.read(1)
            if not byte:
                break
            print("Added to Grid: ", int.from_bytes(byte,  byteorder='little'))
            grid.append(int.from_bytes(byte,  byteorder='little'))
            
        return (rows, columns)
    

        

dimensions = readMazeFile("C:/Users/martm/OneDrive/Desktop/projects/maze to blender/blend.maze")
rows = dimensions[0]
columns = dimensions[1]

def createMaze(grid, columns):

    cell_w = width * 2 + border_width_2
    cell_h = height * 2 + border_width_2

    for i, walls in enumerate(grid):
        # print("walls: ", walls) 
        coords = array_coords_to_matrix_coords(i, columns) 
        row = coords[0] 
        column = coords[1] 
        if(walls & TOP ):
             # +x 
            bpy.ops.mesh.primitive_cube_add(location=(column*(width*2 + border_width_2)+(width) , row*(height*2 + border_width_2) , tall), scale=(width, border_width, tall)) 
            
        if(row == rows-1):
             bpy.ops.mesh.primitive_cube_add(location=(column*(width*2 + border_width_2)+(width) , row*(height*2 + border_width_2)+(height*2 + border_width) , tall), scale=(width, border_width, tall)) 
             
        if(walls & RIGHT): 
            # +y 
            bpy.ops.mesh.primitive_cube_add(location=(column*(width*2 + border_width_2)+(width*2 + border_width), row*(height*2 + border_width_2) + height+border_width, tall), scale=(border_width, height + border_width_2 , tall))
            
        if(column == 0): 
            # +y 
            bpy.ops.mesh.primitive_cube_add(location=(column*(width*2 + border_width_2), row*(height*2 + border_width_2) + height+border_width, tall), scale=(border_width, height + border_width_2 , tall))




        
## +x
#bpy.ops.mesh.primitive_cube_add(location=(1.2,1,1), scale=(0.1, 0.6, 0.5))
## -x
#bpy.ops.mesh.primitive_cube_add(location=(-.2,1,1), scale=(0.1, 0.6, 0.5))
## -y
#bpy.ops.mesh.primitive_cube_add(location=(.5,.5,1), scale=(.6, 0.1, 0.5))
## +y
#bpy.ops.mesh.primitive_cube_add(location=(.5, 1.5,1), scale=(.6, 0.1, 0.5))

createMaze(grid, columns)


