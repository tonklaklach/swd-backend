import logging

from django.http import JsonResponse
from rest_framework.exceptions import APIException, ValidationError
from rest_framework.response import Response
from rest_framework.views import exception_handler as drf_exception_handler

logger = logging.getLogger(__name__)
ERROR_CODES = {
    400: "BAD_REQUEST",
    404: "NOT_FOUND",
    409: "CONFLICT",
    415: "UNSUPPORTED_MEDIA_TYPE",
    422: "VALIDATION_FAILED",
    500: "INTERNAL_ERROR",
    503: "SERVICE_UNAVAILABLE",
}


class Conflict(APIException):
    status_code = 409
    default_detail = "Resource conflicts with existing or dependent data"


class InvalidReference(APIException):
    status_code = 422
    default_detail = "Referenced resource does not exist"
    default_code = "invalid_reference"


def error_body(status, message, code=None):
    return {"code": code or ERROR_CODES.get(status, "HTTP_ERROR"), "error": message}


def exception_handler(exc, context):
    response = drf_exception_handler(exc, context)
    if response is None:
        logger.error("Unexpected API failure", exc_info=exc)
        return Response(error_body(500, "Internal server error"), status=500)

    if isinstance(exc, ValidationError):
        response.status_code = 422
    if "detail" in response.data:
        detail = response.data["detail"]
        if isinstance(detail, list):
            detail = detail[0]
    else:
        field, messages = next(iter(response.data.items()))
        message = messages[0] if isinstance(messages, list) else messages
        if isinstance(message, dict):
            index, nested_messages = next(iter(message.items()))
            field = f"{field}[{index}]"
            message = nested_messages[0]
        detail = f"{field}: {message}"
    code = "INVALID_REFERENCE" if isinstance(exc, InvalidReference) else None
    response.data = error_body(response.status_code, str(detail), code)
    return response


def not_found(request, exception):
    return JsonResponse(error_body(404, "Resource not found"), status=404)
