#ifndef _CppWebImp_H_
#define _CppWebImp_H_

#include "servant/Application.h"

/**
 *
 *
 */
class CppWebImp : public Servant
{
public:
    /**
     *
     */
    virtual ~CppWebImp() {}

    /**
     *
     */
    virtual void initialize();

    /**
     *
     */
    virtual void destroy();

    /**
     *
     */
	int doRequest(TarsCurrentPtr current, vector<char> &buffer);

private:
	int doRequest(const TC_HttpRequest &req, TC_HttpResponse &rsp);
};
/////////////////////////////////////////////////////
#endif
