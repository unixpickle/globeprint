import numpy as np
from PIL import Image

img = Image.open('equi1.jpeg')
arr = np.array(img)
map1 = np.logical_and(arr[:, :, 2:] > arr[:, :, 0:1] * 2,
                      arr[:, :, 2:] > arr[:, :, 1:2]).astype('uint8')
map2 = (1 - map1)
colored = (np.array([[[0x0a, 0xba, 0xb5]]], dtype='uint8') * map1 +
           np.array([[[0xff, 0xff, 0xff]]], dtype='uint8') * map2)
Image.fromarray(colored).save('equi2.png')
