import logging
import requests
import sys
import cv2
import os
import numpy as np
from typing import List, Tuple


def images_to_video(
    name: str,
    images: List[Tuple[str, bytes]]
):
    if not images:
        raise ValueError('requires at least one image')

    images = [(im_name, cv2.imdecode(np.frombuffer(im_data, np.uint8), cv2.IMREAD_COLOR))
              for im_name, im_data in images]

    height, width, *_ = images[0][1].shape

    video_name = f'{name}.avi'
    video = cv2.VideoWriter(video_name, 0, 1/5, (width, height))

    for im_name, im_data in images:
        im_data = cv2.resize(im_data, (width, height))
        cv2.putText(
            img=im_data,
            text=im_name,
            org=(
                int(width*0.1),
                int(height*0.8),
            ),
            fontFace=cv2.FONT_HERSHEY_COMPLEX,
            fontScale=0.5,
            color=(255, 0, 255),
            thickness=1,
        )
        video.write(im_data)

    cv2.destroyAllWindows()
    video.release()

if __name__ == '__main__':
    if len(sys.argv) != 4:
        logging.error('invalid argument count')
        exit(1)

    token = sys.argv[1]
    group_id = sys.argv[2]
    output_file_name = sys.argv[3]

    offset = 0

    videos = list()

    while True:
        resp = requests.post('https://api.vk.com/method/video.get', data={
            'owner_id': group_id,
            'access_token': token,
            'count': 100,
            'offset': offset,
            'v': '5.131'
        })

        items = resp.json()['response']['items']
        if len(items) == 0: 
            break

        videos_info = [{
            'title': item['title'],
            'likes': item['likes']['count'],
            'preview': requests.get(item['image'][-3]['url']).content,
        } for item in items]

        offset += 100

        videos.extend(videos_info)

    videos.sort(key=lambda x: -x['likes'])

    images_to_video(output_file_name, 
    [
        (item['title'], item['preview'])
        for item in videos
    ])



