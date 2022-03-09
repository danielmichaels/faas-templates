import os


def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """

    dirname = os.path.dirname(__file__)
    path = os.path.join(dirname, 'html', 'home.html')

    with open(path, 'r') as file:
        html = file.read()

    return html
