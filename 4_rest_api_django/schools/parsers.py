from io import BytesIO

from rest_framework.exceptions import ParseError
from rest_framework.parsers import JSONParser


class BoundedJSONParser(JSONParser):
    def parse(self, stream, media_type=None, parser_context=None):
        body = stream.read((1 << 20) + 1)
        if len(body) > 1 << 20:
            raise ParseError("Request body must not exceed 1 MiB")
        data = super().parse(BytesIO(body), media_type, parser_context)
        if not isinstance(data, dict):
            raise ParseError("Body must contain one JSON object")
        return data
