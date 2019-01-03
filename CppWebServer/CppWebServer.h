#ifndef _CppWebServer_H_
#define _CppWebServer_H_

#include <iostream>
#include "servant/Application.h"

using namespace tars;

/**
 *
 **/
class CppWebServer : public Application
{
public:
    /**
     *
     **/
    virtual ~CppWebServer() {};

    /**
     *
     **/
    virtual void initialize();

    /**
     *
     **/
    virtual void destroyApp();
};

extern CppWebServer g_app;

////////////////////////////////////////////
#endif
