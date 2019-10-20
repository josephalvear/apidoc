<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:import href="./locales.xsl" />

<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<xsl:template match="/">
    <html lang="{$curr-lang}">
        <head>
            <title><xsl:value-of select="apidoc/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1" />
            <meta name="generator" content="{document('../config.xml')/config/url}" />
            <link rel="stylesheet" type="text/css" href="{$base-url}apidoc.css" />
            <link rel="icon" type="image/png" href="{$icon}" />
            <link rel="license" href="{apidoc/license/@url}" />
            <script src="{$base-url}apidoc.js"></script>
        </head>
        <body>
            <xsl:call-template name="header" />

            <main>
                <div class="content"><xsl:value-of select="/apidoc/content" /></div>

                <xsl:for-each select="apidoc/api">
                <xsl:sort select="path/@path"/>
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </main>

            <footer>
            <div class="wrap">
                <p><xsl:copy-of select="$locale-footer" /></p>
            </div>
            </footer>
        </body>
    </html>
</xsl:template>

<!-- header 界面元素 -->
<xsl:template name="header">
    <header>
        <div class="wrap">
            <h1>
                <img src="{$icon}" />
                <xsl:value-of select="/apidoc/title" />
                <span class="version">&#160;(<xsl:value-of select="/apidoc/@version" />)</span>
            </h1>

            <div class="menus">
                <div class="menu tags-selector">
                    <span><xsl:copy-of select="$locale-tag" /></span>
                    <ul>
                        <xsl:for-each select="apidoc/tag">
                        <li data-tag="{@name}">
                            <label><input type="checkbox" checked="checked" /><xsl:value-of select="@title" /></label>
                        </li>
                        </xsl:for-each>
                    </ul>
                </div>

                <div class="menu methods-selector">
                    <span><xsl:copy-of select="$locale-method" /></span>
                    <ul>
                        <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                        <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                        <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                        <li data-method="{.}">
                            <label><input type="checkbox" checked="true" /><xsl:value-of select="." /></label>
                        </li>
                        </xsl:for-each>
                    </ul>
                </div>

                <div class="menu languages-selector">
                    <span><xsl:copy-of select="$locale-language" /></span>
                    <ul><xsl:call-template name="languages" /></ul>
                </div>
            </div> <!-- end .menus -->
        </div> <!-- end .wrap -->
    </header>
</xsl:template>

<!-- api 界面元素 -->
<xsl:template match="/apidoc/api">
    <xsl:variable name="id" select="concat(@method, translate(path/@path, $id-from, $id-to))" />

    <details id="{$id}" class="api" data-method="{@method}">
    <xsl:attribute name="data-tags">
        <xsl:for-each select="tag"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
    </xsl:attribute>

        <summary>
            <a class="link" href="#{$id}">&#128279;</a> <!-- 链接符号 -->

            <span class="action"><xsl:value-of select="@method" /></span>
            <span>
                <xsl:call-template name="deprecated">
                    <xsl:with-param name="deprecated" select="@deprecated" />
                </xsl:call-template>

                <xsl:value-of select="path/@path" />
            </span>

            <span class="summary"><xsl:value-of select="@summary" /></span>
        </summary>
        <div class="description">
            <xsl:if test="./description"><xsl:value-of select="./description" /></xsl:if>
        </div>

        <div class="body">
            <div class="requests">
                <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                <xsl:for-each select="request">
                    <xsl:call-template name="request">
                        <xsl:with-param name="request" select="." />
                        <xsl:with-param name="path" select="../path" />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
            <div class="responses">
                <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>
                <xsl:for-each select="response">
                    <xsl:call-template name="response">
                        <xsl:with-param name="response" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
        </div>

        <xsl:if test="./callback"><xsl:apply-templates select="./callback" /></xsl:if>
    </details>
</xsl:template>

<!-- 回调内容 -->
<xsl:template match="/apidoc/api/callback">
    <div class="callback" data-method="{./@method}">
        <h3><xsl:copy-of select="$locale-callback" /></h3>
        <div class="description">
            <xsl:value-of select="./@summary" />
            <xsl:if test="./description">
                <br /><xsl:value-of select="./description" />
            </xsl:if>
        </div>

        <div class="body">
            <div class="requests">
                <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                <xsl:for-each select="./request">
                    <xsl:call-template name="request">
                        <xsl:with-param name="request" select="." />
                        <xsl:with-param name="path" select="../path" />
                    </xsl:call-template>
                </xsl:for-each>
            </div>

            <xsl:if test="./response">
                <div class="responses">
                    <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>
                    <xsl:for-each select="./response">
                        <xsl:call-template name="response">
                            <xsl:with-param name="response" select="." />
                        </xsl:call-template>
                    </xsl:for-each>
                </div>
            </xsl:if>
        </div> <!-- end .body -->
    </div> <!-- end .callback -->
</xsl:template>

<!-- api/request 的介面元素 -->
<xsl:template name="request">
<xsl:param name="request" />
<xsl:param name="path" />
<div class="request">
    <xsl:if test="$path/param">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-path-param" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/param" />
        </xsl:call-template>
    </xsl:if>

    <xsl:if test="$path/query">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-query" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/query" />
        </xsl:call-template>
    </xsl:if>
    
    <xsl:if test="$request/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$request/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-body" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$request" />
    </xsl:call-template>
</div>
</xsl:template>

<!-- api/response 的界面 -->
<xsl:template name="response">
    <xsl:param name="response" />

    <h5 class="status"><xsl:value-of select="$response/@status" /></h5>

    <xsl:if test="$response/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$response/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-body" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$response" />
    </xsl:call-template>
</xsl:template>

<!-- path param, path query, header 等的界面 -->
<xsl:template name="param">
    <xsl:param name="title" />
    <xsl:param name="param" />
    <xsl:param name="example" /> <!-- 示例代码 -->

    <div class="param">
        <h4 class="title">&#x27a4;&#160;<xsl:copy-of select="$title" /></h4>
        <table>
            <thead>
                <tr>
                    <th><xsl:copy-of select="$locale-var" /></th>
                    <th><xsl:copy-of select="$locale-type" /></th>
                    <th><xsl:copy-of select="$locale-value" /></th>
                    <th><xsl:copy-of select="$locale-description" /></th>
                </tr>
            </thead>
            <tbody>
                <xsl:call-template name="param-list">
                    <xsl:with-param name="param" select="$param" />
                </xsl:call-template>
            </tbody>
        </table>
    </div>
</xsl:template>

<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="param-list">
    <xsl:param name="param" />
    <xsl:param name="parent" /> <!-- 上一级的名称，嵌套对象时可用 -->

    <xsl:for-each select="$param">
    <tr>
    <xsl:call-template name="deprecated">
        <xsl:with-param name="deprecated" select="@deprecated" />
    </xsl:call-template>

        <th>
            <xsl:if test="$parent"><xsl:value-of select="concat($parent, '.')" /></xsl:if>
            <xsl:value-of select="@name" />
        </th>

        <td>
            <xsl:value-of select="@type" />
            <xsl:if test="@array = 'true'"><xsl:value-of select="'[]'" /></xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="@optional = 'true'"><xsl:value-of select="'O'" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="'R'" /></xsl:otherwise>
            </xsl:choose>

            <xsl:if test="@default">
                <xsl:value-of select="concat(' ', @default)" />
            </xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="description"><xsl:value-of select="description" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
            </xsl:choose>
            <xsl:if test="./enum">
                <p><xsl:copy-of select="$locale-enum" /></p>
                <ul>
                <xsl:for-each select="./enum">
                    <li>
                    <xsl:call-template name="deprecated">
                        <xsl:with-param name="deprecated" select="@deprecated" />
                    </xsl:call-template>

                    <xsl:value-of select="@value" />:<xsl:value-of select="." />
                    </li>
                </xsl:for-each>
                </ul>
            </xsl:if>
        </td>
    </tr>

    <xsl:if test="./param">
        <xsl:call-template name="param-list">
            <xsl:with-param name="param" select="./param" />
            <xsl:with-param name="parent" select="./@name" />
        </xsl:call-template>
    </xsl:if>
    </xsl:for-each>
</xsl:template>

<!--
给指定的元素添加已弃用的标记

该模板会给父元素添加 class 和 title 属性，
所以必须要在父元素的任何子元素之前，否则 chrome 和 safari 可能无法正常解析。
-->
<xsl:template name="deprecated">
    <xsl:param name="deprecated" />

    <xsl:if test="$deprecated">
        <xsl:attribute name="class"><xsl:value-of select="'del'" /></xsl:attribute>
        <xsl:attribute name="title">
            <xsl:value-of select="$deprecated" />
        </xsl:attribute>
    </xsl:if>
</xsl:template>

<!-- 用于将 API 地址转换成合法的 ID 标记 -->
<xsl:variable name="id-from" select="'{}/'" />
<xsl:variable name="id-to" select="'__-'" />

<!-- 根据情况获取相应的图标 -->
<xsl:variable name="icon">
    <xsl:choose>
        <xsl:when test="/apidoc/@logo">
            <xsl:value-of select="/apidoc/@logo" />
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="concat($base-url, '../icon.png')" />
        </xsl:otherwise>
    </xsl:choose>
</xsl:variable>

<!-- 获取相对于当前 xsl 文件的基地址 -->
<xsl:variable name="base-url">
    <xsl:apply-templates select="processing-instruction('xml-stylesheet')" />
</xsl:variable>

<xsl:template match="processing-instruction('xml-stylesheet')">
    <xsl:variable name="v1" select="substring-after(., 'href=&quot;')" />
    <!-- NOTE: 此处假定当前文件叫作 apidoc.xsl，如果不是的话，需要另外处理此代码 -->
    <xsl:variable name="v2" select="substring-before($v1, 'apidoc.xsl&quot;')" />
    <xsl:value-of select="$v2" />
</xsl:template>

</xsl:stylesheet>
